package page

import (
	"blogo/router/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryRoute struct {
	*route.Route
}

func NewCategoryRoute(r *route.Route) *CategoryRoute {
	return &CategoryRoute{Route: r}
}

func (r *CategoryRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "categories", r.TemplateData(gin.H{
		"Content": "This is CATEGORIES page.",
	}))
}
