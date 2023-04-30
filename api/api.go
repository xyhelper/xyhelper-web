package api

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"xyhelper-web/config"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xyhelper/chatgpt-go"
)

// Session
func Session(c *gin.Context) {
	auth := false
	if os.Getenv("AUTH_SECRET_KEY") != "" {
		auth = true
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "",
		"data": gin.H{
			"auth":             auth,
			"model":            "ChatGPTUnofficialProxyAPI",
			"kfurl":            config.Kfurl,
			"fixedBaseURI":     config.BaseURI != "",
			"fixedAccessToken": config.AccessToken != "",
		},
	})
}

// VerifyRequest
type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// Verify
func Verify(c *gin.Context) {
	req := &VerifyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})

		return
	}
	if req.Token == os.Getenv("AUTH_SECRET_KEY") {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Error",
			"message": "Token 错误",
			"data":    nil,
		})
	}
}

// ChatProcessRequest
type ChatProcessRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Optins *struct {
		ConversationId  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
	BaseURI     string `json:"baseURI"`     // 基础URI
	AccessToken string `json:"accessToken"` // 访问令牌
	IsGPT4      bool   `json:"isGPT4"`      // 是否为GPT4
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
	// ctx := c.Request.Context()
	// ctx = gctx.WithCtx(ctx)
	ctx := gctx.New()
	if os.Getenv("AUTH_SECRET_KEY") != "" {
		Authorization := c.GetHeader("Authorization")
		if Authorization != "Bearer "+os.Getenv("AUTH_SECRET_KEY") {
			c.JSON(http.StatusOK, gin.H{
				"status":  "Unauthorized",
				"message": "Token 错误",
				"data":    nil,
			})
			return
		}
	}

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
	BaseURI := req.BaseURI
	if os.Getenv("BASE_URI") != "" {
		BaseURI = os.Getenv("BASE_URI")
	}
	if BaseURI == "" {
		BaseURI = "https://freechat.xyhelper.cn"
	}
	AccessToken := req.AccessToken
	if os.Getenv("ACCESS_TOKEN") != "" {
		AccessToken = os.Getenv("ACCESS_TOKEN")
	}

	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken(AccessToken),
		chatgpt.WithTimeout(time.Duration(config.TimeOutMs*1000*1000)),
		chatgpt.WithBaseURI(BaseURI),
	)
	if req.IsGPT4 {
		cli.SetModel("gpt-4")
	}
	//  如果err不为空,循环获取会话
	stream, err := cli.GetChatStream(req.Prompt, req.Optins.ConversationId, req.Optins.ParentMessageId)
	var message string
	kfurl := "![](" + config.Kfurl + ")"
	for err != nil {

		// 如果返回404，说明会话不存在，重新获取会话
		if err.Error() == "send message failed: 404 Not Found" {
			g.Log().Debug(ctx, "会话不存在，重新获取会话", req)
			req.Optins.ConversationId = ""
			stream, err = cli.GetChatStream(req.Prompt)
			continue
		}

		// 如果返回202，且会话ID为空，且 req.BaseURI 不包含 personalchat ，则重新获取会话
		if err.Error() == "send message failed: 202 Accepted" && req.Optins.ConversationId == "" && !strings.Contains(req.BaseURI, "personalchat") {
			g.Log().Debug(ctx, "共享池新会话分配到未登录账号，重新获取会话", req)
			stream, err = cli.GetChatStream(req.Prompt, req.Optins.ConversationId, req.Optins.ParentMessageId)
			continue
		}

		switch err.Error() {
		// 如果返回429，说明请求过于频繁，等待1秒后重新获取会话
		case "send message failed: 429 Too Many Requests":
			message = "当前请求过于频繁，请稍后再试，或在设置中更换接入点。 联系客服 " + kfurl
		// 如果返回202，提示用户会话登陆中，请稍后再试
		case "send message failed: 202 Accepted":
			message = "当前会话登陆中，请稍后再试，或新建会话"
		default:
			message = err.Error() + "，请稍后再试，或联系客服 " + kfurl
		}

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "Error",
				"message": message,
				"data":    nil,
			})
			return
		}
	}
	res := &ChatProcessResponse{}

	// 使用 Stream 方法向客户端发送 SSE 数据
	c.Stream(func(w io.Writer) bool {
		for text := range stream.Stream {
			// g.DumpWithType(text)
			res.Id = text.MessageID
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
			"apiModel":     "xyhelper",
			"reverseProxy": "https://xyhelper.cn",
			"timeoutMs":    gconv.String(config.TimeOutMs/1000) + "秒",
			"socksProxy":   "-",
			"httpsProxy":   "-",
			"balance":      "-",
			"version":      config.Version,
		},
	})
}
