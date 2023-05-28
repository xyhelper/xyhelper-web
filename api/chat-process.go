package api

import (
	"io"
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
	if config.WeChatServer != "" {
		wxOpenId, err := r.Session.Get("wxOpenId")
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": "请先登录",
				"data":    nil,
			})
		}
		if wxOpenId.String() == "" {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Unauthorized",
				"message": "登陆失效，请重新登陆",
				"data":    nil,
			})
		}
	}
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
		// 延时5秒
		time.Sleep(5 * time.Second)

		conversationReq.ConversationID = ""
		response, err = client.Post(ctx, req.BaseURI+"/backend-api/conversation", conversationReq)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": err.Error(),
				"data":    nil,
			})
		}
		defer response.Close()
	}
	if response.StatusCode != 200 {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "请求失败" + response.Status,
			"data":    nil,
		})
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
		if chatResponse.Message.Author.Role != "assistant" {
			continue
		}
		message := ChatProcessResponse{
			Role:            chatResponse.Message.Author.Role,
			Id:              chatResponse.Message.ID,
			ParentMessageId: req.Optins.ParentMessageId,
			ConversationId:  chatResponse.ConversationID,
			Text:            chatResponse.Message.Content.Parts[0],
		}

		r.Response.Writefln("%s", gjson.New(message).MustToJson())
		r.Response.Flush()
		// g.Log().Debug(ctx, "event.Data", text)
		// g.Log().Debug(ctx, "message", message)
	}
	// 输出 [DONE] 信息
	// r.Response.Writefln("data: %s]\n\n", "[DONE]")
	// r.Response.Flush()

}
