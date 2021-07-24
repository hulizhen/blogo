package page

import (
	"net/http"

	"github.com/hulizhen/blogo/router/route"

	"github.com/gin-gonic/gin"
)

type ArchiveRoute struct {
	*route.Route
}

func NewArchiveRoute(r *route.Route) *ArchiveRoute {
	return &ArchiveRoute{Route: r}
}

func (r *ArchiveRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "archives", r.TemplateData(gin.H{
		"Content": "This is ARCHIVES page.",
	}))
}
