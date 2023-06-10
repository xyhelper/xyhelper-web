package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	TimeOutMs         = 180000
	Version           = "dev"
	Kfurl             = "qrcode_for_gh_19eb5b090f33_258.jpg"
	BaseURI           = ""
	AccessToken       = ""
	ShowPlusBtn       = ""
	AdMessage         = "Aha~"
	PORT              = 8080
	AUTH_SECRET_KEY   = ""
	WeChatServer      = ""
	WechatServerToken = ""
	ShowAbout         = false
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
	wechatServer := g.Cfg().MustGetWithEnv(ctx, "WECHAT_SERVER")
	if !wechatServer.IsEmpty() {
		WeChatServer = wechatServer.String()
	}
	weChatServerToken := g.Cfg().MustGetWithEnv(ctx, "WECHAT_SERVER_TOKEN")
	if !weChatServerToken.IsEmpty() {
		WechatServerToken = weChatServerToken.String()
	}
	showAbout := g.Cfg().MustGetWithEnv(ctx, "SHOW_ABOUT")
	if !showAbout.IsEmpty() {
		ShowAbout = showAbout.Bool()
	}

}
