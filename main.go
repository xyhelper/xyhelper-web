package main

import (
	"embed"
	"fmt"
	"mime"
	"net/http"

	"github.com/gabriel-vasile/mimetype"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

//go:embed frontend/dist
var frontend embed.FS

func main() {
	dirEntries, err := frontend.ReadDir("frontend/dist/assets")
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		fmt.Println(dirEntry.Name())
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 将 / 重定向到 frontend/dist/index.html
	r.GET("/", func(c *gin.Context) {
		content, err := frontend.ReadFile("frontend/dist/index.html")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.Data(http.StatusOK, "text/html", content)
	})
	// 将 /assets/* 重定向到 frontend/dist/assets/*
	r.GET("/assets/*path", func(c *gin.Context) {
		path := c.Param("path")
		ctx := c.Request.Context()
		g.Log().Debug(ctx, "path", path)
		content, err := frontend.ReadFile("frontend/dist/assets" + path)
		if err != nil {
			g.Log().Error(ctx, "err", err)
			c.AbortWithStatus(http.StatusNotFound)
		}
		extension := gfile.ExtName(path)
		extension = "." + extension
		g.Log().Debug(ctx, "extension", extension)
		mimeType := mime.TypeByExtension(extension)
		if mimeType == "" {
			mimeType = mimetype.Detect(content).String()
		}
		c.Data(http.StatusOK, mimeType, content)
	})
	// favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		content, err := frontend.ReadFile("frontend/dist/favicon.ico")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.Data(http.StatusOK, "image/x-icon", content)
	})
	// favion.svg
	r.GET("/favicon.svg", func(c *gin.Context) {
		content, err := frontend.ReadFile("frontend/dist/favicon.svg")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.Data(http.StatusOK, "image/svg+xml", content)
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
