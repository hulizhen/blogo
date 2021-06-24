package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getHome(c *gin.Context) {
	r.store.ArticleStore.ReadArticles()
	c.HTML(http.StatusOK, "home", r.templateData(gin.H{
		"Content": "This is HOME page.",
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
