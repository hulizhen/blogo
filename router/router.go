package router

import (
	"fmt"
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

const staticFilePath = "/dist"

func New(cfg *config.Config, store *store.Store) *Router {
	return &Router{engine: gin.Default(), config: cfg, store: store}
}

func (r *Router) Run() (err error) {
	cfg := r.config
	e := r.engine

	var dist string
	if gin.IsDebugging() {
		dist = "./website"
	} else {
		dist = "./dist"
	}
	e.Static(staticFilePath, dist)
	e.StaticFile("/favicon.ico", cfg.Website.FaviconPath)
	e.StaticFile("/logo.png", cfg.Website.LogoPath)

	e.GET("/", r.getHome)
	e.GET("/archives", r.getArchives)
	e.GET("/categories", r.getCategories)
	e.GET("/tags", r.getTags)
	e.GET("/about", r.getAbout)

	err = r.loadTemplates()
	if err != nil {
		return
	}

	addr := fmt.Sprintf(":%d", r.config.Server.Port)
	r.engine.Run(addr)
	return nil
}

func (r *Router) loadTemplates() (err error) {
	cfg := r.config
	render := multitemplate.NewRenderer()

	p := filepath.Join(cfg.Website.TemplatePath, "include/*.html")
	include, err := filepath.Glob(p)
	if err != nil {
		return
	}

	p = filepath.Join(cfg.Website.TemplatePath, "page/*.html")
	pages, err := filepath.Glob(p)
	if err != nil {
		return
	}

	for _, page := range pages {
		// Make a copy of the 'include' slice to avoid sharing and modifing the same backing array.
		cp := make([]string, len(include))
		copy(cp, include)
		files := append(cp, page)
		name := strings.TrimSuffix(filepath.Base(page), filepath.Ext(page))
		render.AddFromFiles(name, files...)
		// render.AddFromFilesFuncs(name, r.templateFuncMap(), files...)
	}

	r.engine.HTMLRender = render
	return
}

func (r *Router) templateData(data gin.H) gin.H {
	var scriptFilename string
	var styleFilename string
	if gin.IsDebugging() {
		styleFilename = "main.css"
		scriptFilename = "main.js"
	} else {
		styleFilename = "main.min.css"
		scriptFilename = "main.min.js"
	}
	base := gin.H{
		"WebsiteTitle":   r.config.Website.Title,
		"WebsiteAuthor":  r.config.Website.Author,
		"StyleFilename":  filepath.Join(staticFilePath, "style", styleFilename),
		"ScriptFilename": filepath.Join(staticFilePath, "script", scriptFilename),
	}

	for k, v := range base {
		if _, found := data[k]; !found {
			data[k] = v
		}
	}
	return data
}

// func (r *Router) templateFuncMap() template.FuncMap {
// 	fs := []interface{}{
// 		xtime.ShortFormat,
// 		xtime.LongFormat,
// 	}
//
// 	// Iterate the func slice and build a func map of which each key is the last part of func name.
// 	fm := make(template.FuncMap, len(fs))
// 	for _, f := range fs {
// 		ptr := reflect.ValueOf(f).Pointer()
// 		name := runtime.FuncForPC(ptr).Name()
// 		strs := strings.Split(name, ".")
// 		last := strs[len(strs)-1]
// 		fm[last] = f
// 	}
// 	return fm
// }
