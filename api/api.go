package api

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/xyhelper/chatgpt-go"
)

// Session
func Session(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "",
		"data": gin.H{
			"auth":  false,
			"model": "ChatGPTUnofficialProxyAPI",
		},
	})
}

// ChatProcessRequest
type ChatProcessRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Optins *struct {
		ConversationId  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
}

// ChatProcessResponse
type ChatProcessResponse struct {
	Role            string `json:"role"`            // 角色
	Id              string `json:"id"`              // 消息ID
	ParentMessageId string `json:"parentMessageId"` // 父消息ID
	ConversationId  string `json:"conversationId"`  // 会话ID
	Text            string `json:"text"`            // 消息内容
}

// ChatProcess 响应
func ChatProcess(c *gin.Context) {
	var req ChatProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})

		return
	}
	// g.DumpWithType(req)

	token := "xyhelper1231"
	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken(token),
		chatgpt.WithTimeout(180*time.Second),
		chatgpt.WithBaseURI("https://freechat.lidong.xin"),
	)
	stream, err := cli.GetChatStream(req.Prompt, req.Optins.ConversationId, req.Optins.ParentMessageId)
	// 如果返回404，说明会话不存在，重新获取会话
	if err != nil {
		if err.Error() == "send message failed: 404 Not Found" {
			stream, err = cli.GetChatStream(req.Prompt)
		}
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	res := &ChatProcessResponse{}

	// 使用 Stream 方法向客户端发送 SSE 数据
	c.Stream(func(w io.Writer) bool {
		for text := range stream.Stream {
			// g.DumpWithType(text)
			res.Text = text.Content
			res.Role = "assistant"
			res.ConversationId = text.ConversationID
			res.ParentMessageId = req.Optins.ParentMessageId
			data := gjson.New(res).MustToJson()
			writeSSEData(w, data)

		}
		return false
	})

}

func writeSSEData(w io.Writer, data []byte) error {
	_, err := w.Write(append(data, byte('\n')))
	if err != nil {
		return err
	}
	w.(http.Flusher).Flush()

	return nil
}

// Message
func Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "",
		"data": gin.H{
			"apiModel":     "ChatGPTUnofficialProxyAPI",
			"reverseProxy": "https://freechat.lidong.xin/backend-api/conversation",
			"timeoutMs":    100000,
			"socksProxy":   "-",
			"httpsProxy":   "-",
			"balance":      "-",
		},
	})
}
