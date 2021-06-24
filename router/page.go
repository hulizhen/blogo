package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getHome(c *gin.Context) {
	a, err := r.store.ArticleStore.ReadArticles()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "home", r.templateData(gin.H{
		"Articles": a,
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
