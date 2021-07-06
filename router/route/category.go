package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryRoute struct {
	*Route
}

func NewCategoryRoute(r *Route) *CategoryRoute {
	return &CategoryRoute{Route: r}
}

func (r *CategoryRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "categories", r.templateData(gin.H{
		"Content": "This is CATEGORIES page.",
	}))
}
