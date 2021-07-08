package route

import (
	"blogo/config"
	"blogo/store"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Config *config.Config
	Store  *store.Store
}

func (r *Route) LogoPath() string {
	return filepath.Join("/", filepath.Base(r.Config.Website.LogoPath))
}

func (r *Route) TemplateData(data gin.H) gin.H {
	// Copyright year.
	var year string
	since := r.Config.Website.SinceYear
	now := time.Now().Local().Year()
	if since == 0 || since >= now {
		year = strconv.Itoa(now)
	} else {
		year = fmt.Sprintf("%d-%d", since, now)
	}

	base := gin.H{
		"WebsiteTitle":    r.Config.Website.Title,
		"WebsiteAuthor":   r.Config.Website.Author,
		"WebsiteLogoPath": r.LogoPath(),
		"CopyrightYear":   year,
	}

	for k, v := range base {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	return data
}
