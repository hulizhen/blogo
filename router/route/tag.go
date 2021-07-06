package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagRoute struct {
	*Route
}

func NewTagRoute(r *Route) *TagRoute {
	return &TagRoute{Route: r}
}

func (r *TagRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "tags", r.templateData(gin.H{
		"Content": "This is TAGS page.",
	}))
}
