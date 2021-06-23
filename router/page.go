package router

import "github.com/gin-gonic/gin"

func (r *Router) getHome(c *gin.Context) {
	c.HTML(200, "home", gin.H{
		"title":       r.config.Website.Title,
		"description": r.config.Website.Description,
		"content":     "This is HOME page.",
	})
}

func (r *Router) getArchives(c *gin.Context) {
	c.HTML(200, "archives", gin.H{
		"title":       r.config.Website.Title,
		"description": r.config.Website.Description,
		"content":     "This is ARCHIVES page.",
	})
}

func (r *Router) getCategories(c *gin.Context) {
	c.HTML(200, "categories", gin.H{
		"title":       r.config.Website.Title,
		"description": r.config.Website.Description,
		"content":     "This is CATEGORIES page.",
	})
}

func (r *Router) getTags(c *gin.Context) {
	c.HTML(200, "tags", gin.H{
		"title":       r.config.Website.Title,
		"description": r.config.Website.Description,
		"content":     "This is TAGS page.",
	})
}

func (r *Router) getAbout(c *gin.Context) {
	c.HTML(200, "about", gin.H{
		"title":       r.config.Website.Title,
		"description": r.config.Website.Description,
		"content":     "This is ABOUT page.",
	})
}
