package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleRoute struct {
	*Route
}

func NewArticleRoute(r *Route) *ArticleRoute {
	return &ArticleRoute{Route: r}
}

func (r *ArticleRoute) GET(c *gin.Context) {
	slug := c.Param("slug")
	article, err := r.Store.ArticleStore.ReadArticle(slug)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "article", r.templateData(gin.H{
		"Article": article,
	}))
}
