package controller

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xyhelper/chatgpt-go"
	_ "github.com/xyhelper/xyhelper-web/modules/chatweb/controller/admin"
	_ "github.com/xyhelper/xyhelper-web/modules/chatweb/controller/app"
	_ "github.com/xyhelper/xyhelper-web/modules/chatweb/service"
)

type chatreq struct {
	g.Meta  `path:"/chat-process-gf" method:"POST"`
	Prompt  string `json:"prompt"` // 提示
	Options *struct {
		ConversationID  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
}

func init() {
	s := g.Server()
	s.BindHandler("/app/chatweb/api/chat-process", func(r *ghttp.Request) {
		var req *chatreq
		ctx := r.Context()
		err := r.Parse(&req)
		if err != nil {
			r.Response.WriteJson(g.Map{
				"status":  "Error",
				"message": err.Error(),
				"data":    nil,
			})
			r.Exit()
		}
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
			return
		}
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

			r.Response.Write(responseData)
			// r.Response.Write([]byte("\n"))
			r.Response.Flush()
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
			return
		}
		g.Log().Debug(ctx, "写入响应内容完成")
		r.ExitAll()
	})
}
