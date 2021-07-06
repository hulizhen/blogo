package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"blogo/config"
	"blogo/router/route"
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
	e.Static(route.DistFilePath, "./dist")
	e.StaticFile("/favicon.ico", c.Website.FaviconPath)
	e.StaticFile(r.LogoPath(), c.Website.LogoPath)

	// Register routes.
	r.registerRoute("/", route.NewHomeRoute(r.Route))
	r.registerRoute("/archives", route.NewArchiveRoute(r.Route))
	r.registerRoute("/categories", route.NewCategoryRoute(r.Route))
	r.registerRoute("/tags", route.NewTagRoute(r.Route))
	r.registerRoute("/about", route.NewAboutRoute(r.Route))
	r.registerRoute("/articles/:slug", route.NewArticleRoute(r.Route))

	// Load templates.
	err = r.loadTemplates()
	if err != nil {
		return
	}

	// Run on the specified address:port.
	a := fmt.Sprintf(":%d", r.Config.Server.Port)
	r.engine.Run(a)
	return nil
}

// registerRoute register `path` with `route`, which handles the request depending on whether
// the `route` has implemented the cooresponding HTTP method, and aborts with 405 if not.
func (r *Router) registerRoute(path string, route interface{}) {
	r.engine.Any(path, func(c *gin.Context) {
		v := reflect.ValueOf(route)
		m := v.MethodByName(c.Request.Method)

		if !m.IsValid() {
			// The requested HTTP method is not allowed, here we return HTTP status code 405
			// with an "Allow" header field indicating the methods that we have implemented.
			ms := []string{}
			for i := 0; i < len(validHTTPMethods); i++ {
				m = v.MethodByName(validHTTPMethods[i])
				if m.IsValid() {
					ms = append(ms, validHTTPMethods[i])
				}
			}
			c.Header("Allow", strings.Join(ms, ", "))
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		// Call the requested HTTP method implemented by the `route`.
		m.Call([]reflect.Value{reflect.ValueOf(c)})
	})
}

// loadTemplates loads all templates into a map, each value of which is a copy of
// all templates in 'include' directory followed with a single template in 'page' directory.
func (r *Router) loadTemplates() (err error) {
	c := r.Config
	render := multitemplate.NewRenderer()

	p := filepath.Join(c.Website.TemplatePath, "include/*.html")
	include, err := filepath.Glob(p)
	if err != nil {
		return
	}

	p = filepath.Join(c.Website.TemplatePath, "page/*.html")
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
