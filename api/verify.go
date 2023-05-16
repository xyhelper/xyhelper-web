package api

import (
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Verify(r *ghttp.Request) {
	token := r.Get("token").String()
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
