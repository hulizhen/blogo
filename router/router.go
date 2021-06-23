package router

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/hulizhen/blogo/config"
	"gorm.io/gorm"
)

type Router struct {
	engine *gin.Engine
	config *config.Config
	db     *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) *Router {
	return &Router{engine: gin.Default(), config: cfg, db: db}
}

func (r *Router) Run() {
	cfg := r.config
	e := r.engine

	e.Static("/static", "./static")
	e.StaticFile("/favicon.ico", cfg.Website.FaviconPath)
	e.StaticFile("/logo.png", cfg.Website.LogoPath)

	e.GET("/", r.getHome)
	e.GET("/archives", r.getArchives)
	e.GET("/categories", r.getCategories)
	e.GET("/tags", r.getTags)
	e.GET("/about", r.getAbout)

	r.loadTemplates()

	addr := fmt.Sprintf(":%d", r.config.Server.Port)
	r.engine.Run(addr)
}

func (r *Router) loadTemplates() {
	cfg := r.config
	render := multitemplate.NewRenderer()

	p := filepath.Join(cfg.Website.TemplatePath, "base/*.html")
	bases, err := filepath.Glob(p)
	if err != nil {
		log.Panicf("Failed to load templates in '%v'.", p)
	}

	p = filepath.Join(cfg.Website.TemplatePath, "page/*.html")
	pages, err := filepath.Glob(p)
	if err != nil {
		log.Panicf("Failed to load templates in '%v'.", p)
	}

	for _, page := range pages {
		name := strings.TrimSuffix(filepath.Base(page), filepath.Ext(page))
		files := append(bases, page)
		render.AddFromFiles(name, files...)
	}

	r.engine.HTMLRender = render
}
