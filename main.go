package main

import (
	"xyhelper-web/api"
	"xyhelper-web/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	s := g.Server()
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
	s.Run()
}
