package api

import (
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func PushToken(r *ghttp.Request) {
	ctx := r.GetCtx()
	token := r.Get("token").String()
	if !strings.Contains(token, "sup=1|rid=") {
		g.Log().Error(ctx, "token error", token)
		r.Response.WriteJson(g.Map{
			"code": 0,
			"msg":  "token error",
		})
		return
	}
	g.Log().Info(ctx, "token", token)
	res, err := g.Client().ContentJson().Post(ctx, "https://chatarkose.xyhelper.cn/pushtoken", g.Map{
		"token": token,
	})
	if err != nil {
		g.Log().Error(ctx, "pushtoken error", err)
		r.Response.WriteJson(g.Map{
			"code": 0,
			"msg":  "pushtoken error",
		})
		return
	}
	defer res.Close()
	resStr := res.ReadAllString()
	r.Response.WriteJson(gjson.New(resStr))
}
