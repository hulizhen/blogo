package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"blogo/config"
	"blogo/router/route"
	"blogo/router/route/page"
	"blogo/router/route/webhook"
	"blogo/store"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*route.Route
	engine *gin.Engine
}

// Takes from Go source code 'net/http/method.go'.
var validHTTPMethods = [...]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

const distFilePath = "/dist"

// New creates a router instance with an internal router engine.
func New(c *config.Config, s *store.Store) *Router {
	return &Router{
		Route: &route.Route{
			Config: c,
			Store:  s,
		},
		engine: gin.Default(),
	}
}

// Run sets everything up and runs the server.
func (r *Router) Run() (err error) {
	c := r.Config
	e := r.engine

	// Serve static files.
	e.Static(distFilePath, "./dist")
	e.StaticFile("/favicon.ico", c.Website.FaviconPath)
	e.StaticFile(r.LogoPath(), c.Website.LogoPath)

	// Register page routes.
	r.registerRoute("/", page.NewHomeRoute(r.Route))
	r.registerRoute("/archives", page.NewArchiveRoute(r.Route))
	r.registerRoute("/categories", page.NewCategoryRoute(r.Route))
	r.registerRoute("/tags", page.NewTagRoute(r.Route))
	r.registerRoute("/about", page.NewAboutRoute(r.Route))
	r.registerRoute("/articles/:slug", page.NewArticleRoute(r.Route))

	// Register webhook routes.
	r.registerRoute("/webhook/github", webhook.NewGitHubRoute(r.Route))

	// Load templates.
	err = r.loadTemplates()
	if err != nil {
		return
	}

	// Run on the specified address:port.
	a := fmt.Sprintf(":%d", r.Config.Server.Port)
	err = r.engine.Run(a)
	return
}

// registerRoute registers a request handler with the given `path` and `route`,
// which implements supported HTTP methods for the corresponding resource.
func (r *Router) registerRoute(path string, route interface{}) {
	v := reflect.ValueOf(route)
	for _, n := range validHTTPMethods {
		m := v.MethodByName(n)
		if m.IsValid() {
			if f, ok := m.Interface().(func(*gin.Context)); ok {
				r.engine.Handle(n, path, f)
			}
		}
	}
}

// loadTemplates loads all templates into a map, each value of which is a copy of
// all templates in 'include' directory followed with a single template in 'page' directory.
func (r *Router) loadTemplates() (err error) {
	c := r.Config
	render := multitemplate.NewRenderer()

	p := filepath.Join(c.Website.TemplatePath, "include/*.gohtml")
	include, err := filepath.Glob(p)
	if err != nil {
		return
	}

	p = filepath.Join(c.Website.TemplatePath, "page/*.gohtml")
	pages, err := filepath.Glob(p)
	if err != nil {
		return
	}

	for _, pg := range pages {
		// Make a copy of the 'include' slice to avoid sharing and modifing the same backing array.
		cp := make([]string, len(include))
		copy(cp, include)
		files := append(cp, pg)
		name := strings.TrimSuffix(filepath.Base(pg), filepath.Ext(pg))
		render.AddFromFiles(name, files...)
		// render.AddFromFilesFuncs(name, r.templateFuncMap(), files...)
	}

	r.engine.HTMLRender = render
	return
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
