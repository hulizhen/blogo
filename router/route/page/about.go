package page

import (
	"github.com/hulizhen/blogo/router/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AboutRoute struct {
	*route.Route
}

func NewAboutRoute(r *route.Route) *AboutRoute {
	return &AboutRoute{Route: r}
}

func (r *AboutRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "about", r.TemplateData(gin.H{
		"About": r.Store.AboutStore.About,
	}))
}
