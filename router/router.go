package router

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"blogo/config"
	"blogo/store"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	config *config.Config
	store  *store.Store
}

func New(cfg *config.Config, store *store.Store) *Router {
	return &Router{engine: gin.Default(), config: cfg, store: store}
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
		// Make a copy of the 'bases' slice to avoid sharing and modifing the same backing array.
		cp := make([]string, len(bases))
		copy(cp, bases)
		files := append(cp, page)
		name := strings.TrimSuffix(filepath.Base(page), filepath.Ext(page))
		render.AddFromFiles(name, files...)
	}

	r.engine.HTMLRender = render
}

func (r *Router) templateData(data gin.H) gin.H {
	base := gin.H{
		"WebsiteTitle":  r.config.Website.Title,
		"WebsiteAuthor": r.config.Website.Author,
	}
	for k, v := range base {
		if _, found := data[k]; !found {
			data[k] = v
		}
	}
	return data
}
