package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	TimeOutMs = 180000
	Version   = "dev"
	Kfurl     = "https://work.weixin.qq.com/kf/kefu/qrcode?kfcode=kfc97c97206f588c396"
)

func init() {
	// TimeOutMs = 180000
	ctx := gctx.GetInitCtx()
	timeOutMs := g.Cfg().MustGetWithEnv(ctx, "TIMEOUTMS")
	if !timeOutMs.IsEmpty() {
		TimeOutMs = timeOutMs.Int()
	}
	kfurl := g.Cfg().MustGetWithEnv(ctx, "KFURL")
	if !kfurl.IsEmpty() {
		Kfurl = kfurl.String()
	}

}
