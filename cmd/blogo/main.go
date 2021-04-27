package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("../../template/*")
	engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":   "Blogo",
			"content": "A blog engine built with Go.",
		})
	})
	engine.Run(":8080")
}
