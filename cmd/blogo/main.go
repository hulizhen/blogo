package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hulizhen/blogo/config"
	"github.com/hulizhen/blogo/service/observer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	cfg := config.New(config.ConfigFilePath)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%d)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Panicf("Failed to open database with error: %v.", err)
	}

	// Start observing the repo changes.
	observer, err := observer.NewRepoObserver(db, cfg.Website.BlogRepoPath)
	if err != nil {
		log.Panicf("Failed to create repo observer with error: %v.", err)
	}
	observer.Start()

	e := gin.Default()
	e.LoadHTMLGlob("template/*")
	e.Static("/static", "./static")
	e.StaticFile("/favicon.ico", cfg.Website.FaviconPath)
	e.StaticFile("/logo.png", cfg.Website.LogoPath)
	e.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":       cfg.Website.Title,
			"description": cfg.Website.Description,
		})
	})
	e.GET("/archives", func(c *gin.Context) {
		c.HTML(200, "archives.html", gin.H{
			"title":       cfg.Website.Title,
			"description": cfg.Website.Description,
			"content":     "This is ARCHIVES page.",
		})
	})
	e.GET("/categories", func(c *gin.Context) {
		c.HTML(200, "categories.html", gin.H{
			"title":       cfg.Website.Title,
			"description": cfg.Website.Description,
			"content":     "This is CATEGORIES page.",
		})
	})
	e.GET("/tags", func(c *gin.Context) {
		c.HTML(200, "tags.html", gin.H{
			"title":       cfg.Website.Title,
			"description": cfg.Website.Description,
			"content":     "This is TAGS page.",
		})
	})
	e.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"title":       cfg.Website.Title,
			"description": cfg.Website.Description,
			"content":     "This is ABOUT page.",
		})
	})
	e.Run(addr)
}
