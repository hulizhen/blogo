package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hulizhen/blogo/config"
)

func main() {
	cfg := config.SharedConfig()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	e := gin.Default()
	e.LoadHTMLGlob("template/*")
	e.StaticFile("favicon.ico", cfg.Website.FaviconPath)
	e.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":   cfg.Website.Title,
			"content": cfg.Website.Description,
		})
	})
	e.Run(addr)
}
