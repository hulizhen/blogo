package webhook

import (
	"blogo/router/route"
	"fmt"
	"net/http"
	"os/exec"
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
		c.AbortWithStatus(http.StatusBadRequest)
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
	if url == r.Config.Website.BlogRepoURL && branch == r.Config.Website.BlogRepoBranch {
		cmd := fmt.Sprintf("cd %v && git pull origin %v", r.Config.Website.BlogRepoPath, r.Config.Website.BlogRepoBranch)
		err := exec.Command("/bin/sh", "-c", cmd).Run()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = r.Store.ArticleStore.ScanArticles()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}
