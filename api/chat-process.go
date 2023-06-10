package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
	chatresponse "xyhelper-web/chat-response"
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"github.com/launchdarkly/eventsource"
)

// ChatProcessRequest
type ChatProcessRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Optins *struct {
		ConversationId  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
	BaseURI     string `json:"baseURI"`     // 基础URI
	AccessToken string `json:"accessToken"` // 访问令牌
	// IsGPT4      bool   `json:"isGPT4"`      // 是否为GPT4
	Model string `json:"model"` // 模型
}

// ChatProcessResponse
type ChatProcessResponse struct {
	Role            string `json:"role"`            // 角色
	Id              string `json:"id"`              // 消息ID
	ParentMessageId string `json:"parentMessageId"` // 父消息ID
	ConversationId  string `json:"conversationId"`  // 会话ID
	Text            string `json:"text"`            // 消息内容
	FinishType      string `json:"finishType"`      // 结束类型
}

type Message struct {
	ID      string  `json:"id"`
	Author  Author  `json:"author"`
	Content Content `json:"content"`
}

type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ConversationReq struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	ConversationID  string    `json:"conversation_id,omitempty"`
	ParentMessageID string    `json:"parent_message_id"`
	Model           string    `json:"model"`
	TimezoneOffset  int       `json:"timezone_offset_min"`
	// history_and_training_disabled
	HistoryAndTrainingDisabled bool `json:"history_and_training_disabled"`
	// supports_modapi
	SupportsModapi bool `json:"supports_modapi"`
	// variant_purpose
	// VariantPurpose string `json:"variant_purpose"`
}

func ChatProcess(r *ghttp.Request) {
	ctx := r.Context()
	// if config.WeChatServer != "" {
	// 	wxOpenId, err := r.Session.Get("wxOpenId")
	// 	if err != nil {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"status":  "Error",
	// 			"message": "请先登录",
	// 			"data":    nil,
	// 		})
	// 	}
	// 	if wxOpenId.String() == "" {
	// 		r.Response.WriteJsonExit(g.Map{
	// 			"status":  "Unauthorized",
	// 			"message": "登陆失效，请重新登陆",
	// 			"data":    nil,
	// 		})
	// 	}
	// }
	var req *ChatProcessRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "参数错误",
			"data":    nil,
		})
	}
	// usermodel := "text-davinci-002-render-sha"
	// if req.IsGPT4 {
	// 	usermodel = "gpt-4"
	// }
	parentMessageId := uuid.New().String()
	if req.Optins.ParentMessageId != "" {
		parentMessageId = req.Optins.ParentMessageId
	}
	conversationReq := ConversationReq{
		Action: "next",
		Messages: []Message{
			{
				ID:     uuid.New().String(),
				Author: Author{Role: "user"},
				Content: Content{
					ContentType: "text",
					Parts:       []string{req.Prompt},
				},
			},
		},
		ConversationID:  req.Optins.ConversationId,
		ParentMessageID: parentMessageId,
		Model:           req.Model,
	}
	// g.Dump(conversationReq)

	client := g.Client()
	client.SetHeader("Authorization", "Bearer "+req.AccessToken)
	client.SetHeader("Content-Type", "application/json")
	baseURI := req.BaseURI
	if config.BaseURI != "" {
		baseURI = config.BaseURI
	}
	if config.AccessToken != "" {
		client.SetHeader("Authorization", "Bearer "+config.AccessToken)
	}

	response, err := client.Post(ctx, baseURI+"/backend-api/conversation", conversationReq)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	defer response.Close()
	// 如果返回404，去掉会话ID，重新请求
	if response.StatusCode == 404 {
		g.Log().Error(ctx, "会话ID不存在，重新请求")
		// 延时5秒
		time.Sleep(5 * time.Second)

		conversationReq.ConversationID = ""
		response, err = client.Post(ctx, baseURI+"/backend-api/conversation", conversationReq)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": err.Error(),
				"data":    nil,
			})
		}
		defer response.Close()
	}
	if response.StatusCode == 501 {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "当前用户不支持使用该模型,请联系客服购买会员",
			"data":    nil,
		})
	}
	if response.StatusCode == 429 {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "请求失败" + response.Status + ", 当前请求过多，请稍后再试,或新建聊天窗口",
			"data":    nil,
		})
	}
	if response.StatusCode == 413 {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "请求失败" + response.Status + ", 请求内容过长，请重新输入",
			"data":    nil,
		})
	}

	if response.StatusCode != 200 {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "请求失败" + response.Status + ", 请联系客服",
			"data":    nil,
		})
	}

	//  流式回应
	rw := r.Response.RawWriter()
	flusher, ok := rw.(http.Flusher)
	if !ok {
		g.Log().Error(ctx, "rw.(http.Flusher) error")
		r.Response.WriteStatusExit(500)
		return
	}
	// r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	// r.Response.Flush()
	decoder := eventsource.NewDecoder(response.Body)
	// var message ChatProcessResponse
	for {
		event, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			g.Log().Error(ctx, "decoder.Decode error", err)
			continue
		}
		text := event.Data()
		if text == "" {
			continue
		}
		var chatResponse chatresponse.Conversation
		err = gconv.Struct(text, &chatResponse)
		if err != nil {
			continue
		}
		// g.Log().Debug(ctx, "chatResponse", chatResponse)
		if chatResponse.Message.Author.Role != "assistant" {
			continue
		}
		message := ChatProcessResponse{
			Role:            chatResponse.Message.Author.Role,
			Id:              chatResponse.Message.ID,
			ParentMessageId: req.Optins.ParentMessageId,
			ConversationId:  chatResponse.ConversationID,
			Text:            chatResponse.Message.Content.Parts[0],
			FinishType:      chatResponse.Message.Metadata.FinishDetails.Type,
		}
		// g.Log().Debug(ctx, "message", message)

		_, err = fmt.Fprintf(rw, "%s\n", gjson.New(message).String())
		if err != nil {
			g.Log().Error(ctx, "fmt.Fprintf error", err)
			response.Body.Close()
			continue
		}
		flusher.Flush()
	}

}
