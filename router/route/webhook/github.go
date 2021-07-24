package webhook

import (
	"github.com/hulizhen/blogo/router/route"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GitHubRoute struct {
	*route.Route
}

func NewGitHubRoute(r *route.Route) *GitHubRoute {
	return &GitHubRoute{Route: r}
}

func (r *GitHubRoute) POST(c *gin.Context) {
	payload := make(gin.H)
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	ref, ok := payload["ref"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	url, ok := payload["repository"].(map[string]interface{})["clone_url"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	branch := strings.TrimPrefix(ref, "refs/heads/")

	// Update local repo.
	if url == r.Config.Repository.RemoteURL && branch == r.Config.Repository.Branch {
		err = r.Store.ArticleStore.ScanArticles()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
}
