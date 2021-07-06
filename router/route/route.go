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

const DistFilePath = "/dist"

func (r *Route) LogoPath() string {
	return filepath.Join("/", filepath.Base(r.Config.Website.LogoPath))
}

func styleFilePath() string {
	var filename string
	if gin.IsDebugging() {
		filename = "bundle.css"
	} else {
		filename = "bundle.min.css"
	}
	return filepath.Join(DistFilePath, "style", filename)
}

func scriptFilePaths() []string {
	var filenames []string
	if gin.IsDebugging() {
		filenames = []string{
			"main.js",
		}
	} else {
		filenames = []string{
			"bundle.min.js",
		}
	}
	var filePaths []string
	for _, filename := range filenames {
		filePaths = append(filePaths, filepath.Join(DistFilePath, "script", filename))
	}
	return filePaths
}

func (r *Route) templateData(data gin.H) gin.H {
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
		"StyleFilePath":   styleFilePath(),
		"ScriptFilePaths": scriptFilePaths(),
	}

	for k, v := range base {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	return data
}
