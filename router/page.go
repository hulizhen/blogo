package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"title":   r.config.Website.Title,
		"content": "This is HOME page.",
	})
}

func (r *Router) getArchives(c *gin.Context) {
	c.HTML(http.StatusOK, "archives", gin.H{
		"title":   r.config.Website.Title,
		"content": "This is ARCHIVES page.",
	})
}

func (r *Router) getCategories(c *gin.Context) {
	c.HTML(http.StatusOK, "categories", gin.H{
		"title":   r.config.Website.Title,
		"content": "This is CATEGORIES page.",
	})
}

func (r *Router) getTags(c *gin.Context) {
	c.HTML(http.StatusOK, "tags", gin.H{
		"title":   r.config.Website.Title,
		"content": "This is TAGS page.",
	})
}

func (r *Router) getAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"title":   r.config.Website.Title,
		"content": "This is ABOUT page.",
	})
}
