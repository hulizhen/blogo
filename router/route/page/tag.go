package page

import (
	"github.com/hulizhen/blogo/router/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagRoute struct {
	*route.Route
}

func NewTagRoute(r *route.Route) *TagRoute {
	return &TagRoute{Route: r}
}

func (r *TagRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "tags", r.TemplateData(gin.H{
		"Content": "This is TAGS page.",
	}))
}
