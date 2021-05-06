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
	e.Static("/static", "./static")
	e.StaticFile("/favicon.ico", cfg.Website.FaviconPath)
	e.StaticFile("/logo.png", cfg.Website.LogoPath)
	e.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"name":        cfg.Website.Name,
			"description": cfg.Website.Description,
		})
	})
	e.GET("/archives", func(c *gin.Context) {
		c.HTML(200, "archives.html", gin.H{
			"name":        cfg.Website.Name,
			"description": cfg.Website.Description,
			"content":     "This is ARCHIVES page.",
		})
	})
	e.GET("/categories", func(c *gin.Context) {
		c.HTML(200, "categories.html", gin.H{
			"name":        cfg.Website.Name,
			"description": cfg.Website.Description,
			"content":     "This is CATEGORIES page.",
		})
	})
	e.GET("/tags", func(c *gin.Context) {
		c.HTML(200, "tags.html", gin.H{
			"name":        cfg.Website.Name,
			"description": cfg.Website.Description,
			"content":     "This is TAGS page.",
		})
	})
	e.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"name":        cfg.Website.Name,
			"description": cfg.Website.Description,
			"content":     "This is ABOUT page.",
		})
	})
	e.Run(addr)
}
