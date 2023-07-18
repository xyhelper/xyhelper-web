package main

import (
	"time"
	"xyhelper-web/api"
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gsession"
)

func main() {
	s := g.Server()
	if !gfile.Exists("./data/sessions") {
		gfile.Mkdir("./data/sessions")
	}
	s.SetSessionStorage(gsession.NewStorageFile("./data/sessions", 3600*24*180*time.Second))
	s.SetSessionCookieMaxAge(3600 * 24 * 180 * time.Second)
	s.SetPort(config.PORT)
	if gfile.Exists("frontend/dist") {
		s.SetServerRoot("frontend/dist")
	} else if gfile.Exists("dist") {
		s.SetServerRoot("dist")
	}
	apiGroup := s.Group("/api")
	apiGroup.POST("/session", api.Session)
	apiGroup.POST("/chat-process", api.ChatProcess)
	// apiGroup.POST("/config", api.Config)
	apiGroup.POST("/verify", api.Verify)
	apiGroup.POST("/pushtoken", api.PushToken)
	s.Run()
}
