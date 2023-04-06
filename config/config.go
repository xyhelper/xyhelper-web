package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	TimeOutMs = 180000
	Version   = "dev"
)

func init() {
	// TimeOutMs = 180000
	ctx := gctx.GetInitCtx()
	timeOutMs := g.Cfg().MustGetWithEnv(ctx, "TIMEOUTMS")
	if !timeOutMs.IsEmpty() {
		TimeOutMs = timeOutMs.Int()
	}
}
