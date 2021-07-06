package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArchiveRoute struct {
	*Route
}

func NewArchiveRoute(r *Route) *ArchiveRoute {
	return &ArchiveRoute{Route: r}
}

func (r *ArchiveRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "archives", r.templateData(gin.H{
		"Content": "This is ARCHIVES page.",
	}))
}
