package router

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

const distFilePath = "/dist"

func New(cfg *config.Config, store *store.Store) *Router {
	return &Router{engine: gin.Default(), config: cfg, store: store}
}

func (r *Router) Run() (err error) {
	cfg := r.config
	e := r.engine

	e.Static(distFilePath, "./dist")
	e.StaticFile("/favicon.ico", cfg.Website.FaviconPath)
	e.StaticFile(logoPath(cfg), cfg.Website.LogoPath)

	e.GET("/", r.getHome)
	e.GET("/archives", r.getArchives)
	e.GET("/categories", r.getCategories)
	e.GET("/tags", r.getTags)
	e.GET("/about", r.getAbout)
	e.GET("/article/:slug", r.getArticle)

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

func styleFilePath() string {
	var filename string
	if gin.IsDebugging() {
		filename = "bundle.css"
	} else {
		filename = "bundle.min.css"
	}
	return filepath.Join(distFilePath, "style", filename)
}

func scriptFilePaths() []string {
	var filenames []string
	if gin.IsDebugging() {
		filenames = []string{
			"main.js",
			"prism/prism.js",
		}
	} else {
		filenames = []string{
			"bundle.min.js",
		}
	}
	var filePaths []string
	for _, filename := range filenames {
		filePaths = append(filePaths, filepath.Join(distFilePath, "script", filename))
	}
	return filePaths
}

func (r *Router) templateData(data gin.H) gin.H {
	// Copyright year.
	var year string
	since := r.config.Website.SinceYear
	now := time.Now().Local().Year()
	if since == 0 || since >= now {
		year = strconv.Itoa(now)
	} else {
		year = fmt.Sprintf("%d-%d", since, now)
	}

	base := gin.H{
		"WebsiteTitle":    r.config.Website.Title,
		"WebsiteAuthor":   r.config.Website.Author,
		"WebsiteLogoPath": logoPath(r.config),
		"CopyrightYear":   year,
		"StyleFilePath":   styleFilePath(),
		"ScriptFilePaths": scriptFilePaths(),
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

func logoPath(cfg *config.Config) string {
	return filepath.Join("/", filepath.Base(cfg.Website.LogoPath))
}
