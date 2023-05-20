package api

import (
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Verify(r *ghttp.Request) {
	ctx := r.Context()
	token := r.Get("token").String()
	if config.WeChatServer != "" {
		client := g.Client()
		client.SetHeader("Content-Type", "application/json")
		client.SetHeader("Authorization", "Bearer "+config.WechatServerToken)
		resp, err := client.Get(ctx, config.WeChatServer+"/api/wechat/user?code="+token)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": err.Error(),
				"data":    nil,
			})
		}
		defer resp.Close()
		// g.Dump(resp.ReadAllString())

		if resp.StatusCode != 200 {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": "验证错误," + resp.Status,
				"data":    nil,
			})
		}
		respJson := gjson.New(resp.ReadAllString())
		g.Dump(respJson)
		wxOpenId := respJson.Get("data").String()
		if wxOpenId == "" {
			r.Response.WriteJsonExit(g.Map{
				"status":  "Error",
				"message": "验证码错误",
				"data":    nil,
			})
		}
		r.Session.Set("wxOpenId", wxOpenId)
		r.Response.WriteJsonExit(g.Map{
			"status":  "Success",
			"message": "",
			"data":    nil,
		})
		return
	}
	if token == config.AUTH_SECRET_KEY {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Success",
			"message": "",
			"data":    nil,
		})
	} else {
		r.Response.WriteJsonExit(g.Map{
			"status":  "Error",
			"message": "授权码错误",
			"data":    nil,
		})
	}

}
