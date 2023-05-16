package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	TimeOutMs       = 180000
	Version         = "dev"
	Kfurl           = "https://work.weixin.qq.com/kf/kefu/qrcode?kfcode=kfc97c97206f588c396"
	BaseURI         = ""
	AccessToken     = ""
	ShowPlusBtn     = ""
	AdMessage       = "Aha~"
	PORT            = 8080
	AUTH_SECRET_KEY = ""
)

func init() {
	// TimeOutMs = 180000
	ctx := gctx.GetInitCtx()
	port := g.Cfg().MustGetWithEnv(ctx, "PORT")
	if !port.IsEmpty() {
		PORT = port.Int()
	}
	timeOutMs := g.Cfg().MustGetWithEnv(ctx, "TIMEOUTMS")
	if !timeOutMs.IsEmpty() {
		TimeOutMs = timeOutMs.Int()
	}
	kfurl := g.Cfg().MustGetWithEnv(ctx, "KFURL")
	if !kfurl.IsEmpty() {
		Kfurl = kfurl.String()
	}
	baseURI := g.Cfg().MustGetWithEnv(ctx, "BASE_URI")
	if !baseURI.IsEmpty() {
		BaseURI = baseURI.String()
	}
	accessToken := g.Cfg().MustGetWithEnv(ctx, "ACCESS_TOKEN")
	if !accessToken.IsEmpty() {
		AccessToken = accessToken.String()
	}
	showPlusBtn := g.Cfg().MustGetWithEnv(ctx, "SHOW_PLUS_BTN")
	if !showPlusBtn.IsEmpty() {
		ShowPlusBtn = showPlusBtn.String()
	}
	adMesage := g.Cfg().MustGetWithEnv(ctx, "AD_MESSAGE")
	if !adMesage.IsEmpty() {
		AdMessage = adMesage.String()
	}
	authSecretKey := g.Cfg().MustGetWithEnv(ctx, "AUTH_SECRET_KEY")
	if !authSecretKey.IsEmpty() {
		AUTH_SECRET_KEY = authSecretKey.String()
	}

}
