package route

import (
	"blogo/pkg/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomeRoute struct {
	*Route
}

func NewHomeRoute(r *Route) *HomeRoute {
	return &HomeRoute{Route: r}
}

func (r *HomeRoute) GET(c *gin.Context) {
	// Get offset.
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 1 {
		offset = 1
	}

	// Get articles.
	pageSize := r.Config.Website.ArticlePageSize
	articles, err := r.Store.ArticleStore.ReadArticles(pageSize, offset-1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	// Get article count.
	count, err := r.Store.ArticleStore.ReadArticleCount()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	// Generate pagination.
	pagination := pagination.New(count, pageSize, offset, c.FullPath())

	c.HTML(http.StatusOK, "home", r.templateData(gin.H{
		"Articles":   articles,
		"Pagination": pagination,
	}))
}
