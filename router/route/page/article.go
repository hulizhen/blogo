package page

import (
	"blogo/router/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleRoute struct {
	*route.Route
}

func NewArticleRoute(r *route.Route) *ArticleRoute {
	return &ArticleRoute{Route: r}
}

func (r *ArticleRoute) GET(c *gin.Context) {
	slug := c.Param("slug")
	article, err := r.Store.ArticleStore.ReadArticle(slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "article", r.TemplateData(gin.H{
		"Article": article,
	}))
}
