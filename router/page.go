package router

import (
	"blogo/pkg/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) getHome(c *gin.Context) {
	// Get offset.
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 1 {
		offset = 1
	}

	// Get articles.
	pageSize := pagination.DefaultPageSize
	articles, err := r.store.ArticleStore.ReadArticles(pageSize, offset-1)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	// Get article count.
	count, err := r.store.ReadArticleCount()
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

func (r *Router) getArchives(c *gin.Context) {
	c.HTML(http.StatusOK, "archives", r.templateData(gin.H{
		"Content": "This is ARCHIVES page.",
	}))
}

func (r *Router) getCategories(c *gin.Context) {
	c.HTML(http.StatusOK, "categories", r.templateData(gin.H{
		"Content": "This is CATEGORIES page.",
	}))
}

func (r *Router) getTags(c *gin.Context) {
	c.HTML(http.StatusOK, "tags", r.templateData(gin.H{
		"Content": "This is TAGS page.",
	}))
}

func (r *Router) getAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about", r.templateData(gin.H{
		"Content": "This is ABOUT page.",
	}))
}
