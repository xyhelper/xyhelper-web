package api

import (
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Session(r *ghttp.Request) {
	auth := false
	if config.AUTH_SECRET_KEY != "" {
		auth = true
	}
	if config.WeChatServer != "" {
		auth = true
	}

	r.Response.WriteJsonExit(g.Map{
		"status":  "Success",
		"message": "",
		"data": g.Map{
			"auth":             auth,
			"model":            "ChatGPTUnofficialProxyAPI",
			"kfurl":            config.Kfurl,
			"fixedBaseURI":     config.BaseURI != "",
			"fixedAccessToken": config.AccessToken != "",
			"showPlusBtn":      config.ShowPlusBtn,
			"adMessage":        config.AdMessage,
			"showAbout":        config.ShowAbout,
		},
	})
}
