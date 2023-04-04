package app

import (
	"context"
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	chatgpt "github.com/xyhelper/chatgpt-go"
)

type ChatwebApiController struct {
	*cool.ControllerSimple
}

func init() {
	var chatweb_api_controller = &ChatwebApiController{
		&cool.ControllerSimple{
			Perfix: "/app/chatweb/api",
		},
	}
	// 注册路由
	cool.RegisterControllerSimple(chatweb_api_controller)
}

// 增加 Welcome 演示 方法
type ChatwebApiWelcomeReq struct {
	g.Meta `path:"/welcome" method:"GET"`
}
type ChatwebApiWelcomeRes struct {
	*cool.BaseRes
	Data interface{} `json:"data"`
}

func (c *ChatwebApiController) Welcome(ctx context.Context, req *ChatwebApiWelcomeReq) (res *ChatwebApiWelcomeRes, err error) {
	res = &ChatwebApiWelcomeRes{
		BaseRes: cool.Ok("Welcome to Cool Admin Go"),
		Data:    gjson.New(`{"name": "Cool Admin Go", "age":0}`),
	}
	return
}

// SessionReq 请求参数
type ChatwebApiSessionReq struct {
	g.Meta `path:"/session" method:"POST"`
}

// SessionRes 返回参数
type ChatwebApiSessionRes struct {
	Status  string `json:"status"`  // 状态
	Message string `json:"message"` // 消息
	Data    *struct {
		Auth  bool   `json:"auth"`  // 认证
		Model string `json:"model"` // 模型
	} `json:"data"` // 数据
}

// Session 会话
func (c *ChatwebApiController) Session(ctx context.Context, req *ChatwebApiSessionReq) (res *ChatwebApiSessionRes, err error) {
	res = &ChatwebApiSessionRes{
		Status:  "Success",
		Message: "",
		Data: &struct {
			Auth  bool   `json:"auth"`
			Model string `json:"model"`
		}{
			Auth:  false,
			Model: "ChatGPTUnofficialProxyAPI",
		},
	}
	return
}

// ChatProcessReq 请求参数
type ChatwebApiChatProcessReq struct {
	g.Meta  `path:"/chat-process-gf" method:"POST"`
	Prompt  string `json:"prompt"` // 提示
	Options *struct {
		ConversationID  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
}

// ChatProcessRes 返回参数
type ChatwebApiChatProcessRes struct {
}

// ChatProcess 会话
func (c *ChatwebApiController) ChatProcess(ctx context.Context, req *ChatwebApiChatProcessReq) (res *ChatwebApiChatProcessRes, err error) {
	r := g.RequestFromCtx(ctx)
	token := `random token`

	cli := chatgpt.NewClient(
		// chatgpt.WithDebug(true),
		chatgpt.WithTimeout(120*time.Second),
		chatgpt.WithAccessToken(token),
		chatgpt.WithBaseURI("https://freechat.lidong.xin"),
	)
	var stream *chatgpt.ChatStream
	if req.Options.ConversationID == "" && req.Options.ParentMessageId == "" {
		stream, err = cli.GetChatStream(req.Prompt)
	} else {
		stream, err = cli.GetChatStream(req.Prompt, req.Options.ConversationID, req.Options.ParentMessageId)
	}
	if err != nil {
		return nil, err
	}
	// w := r.Response.RawWriter()
	// f, ok := w.(http.Flusher)
	// if !ok {
	// 	return nil, err
	// }
	g.Log().Debug(ctx, "开始写入响应内容")
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Flush()
	for text := range stream.Stream {
		// msg := gjson.New(text.String())
		responseData := map[string]interface{}{
			"role":            "assistant",
			"id":              text.MessageID,
			"parentMessageId": req.Options.ParentMessageId,
			"conversationId":  text.ConversationID,
			"text":            text.Content,
			// "time":            time.Now().Unix(),
		}
		// fmt.Fprintf(w, "%s\n", gjson.New(responseData).MustToJson())
		// r.Response.ResponseWriter.Write()

		// r.Response.RawWriter().Write(gconv.Bytes(responseData))
		// r.Response.RawWriter().Write([]byte("\n"))
		// r.Response.SetBuffer(gconv.Bytes(responseData))
		r.Response.WriteJson(responseData)

		r.Response.Writer.Flush()

		// resJson := gjson.New(responseData).MustToJson()
		// r.Response.WriteJson(responseData)
		// r.Response.Write([]byte("\n"))
		// r.Response.Flush()
		// r.Response.RawWriter().Write([]byte(resJson))
		// r.Response.Flush()
		g.Log().Debug(ctx, responseData)
	}
	if err := stream.Err; err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	g.Log().Debug(ctx, "写入响应内容完成")
	// r.Exit()

	// res = &ChatwebApiChatProcessRes{}
	return
}
