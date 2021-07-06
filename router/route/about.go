package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AboutRoute struct {
	*Route
}

func NewAboutRoute(r *Route) *AboutRoute {
	return &AboutRoute{Route: r}
}

func (r *AboutRoute) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "about", r.templateData(gin.H{
		"About": r.Store.AboutStore.About,
	}))
}
