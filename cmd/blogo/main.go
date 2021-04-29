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
	e.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":   cfg.Blog.Title,
			"content": cfg.Blog.Description,
		})
	})
	e.Run(addr)
}
