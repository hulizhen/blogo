package observer

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hulizhen/blogo/model"
	"gorm.io/gorm"
)

type RepoObserver struct {
	db       *gorm.DB
	repoPath string
}

func NewRepoObserver(db *gorm.DB, repoPath string) (*RepoObserver, error) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, err
	}

	return &RepoObserver{
		db:       db,
		repoPath: repoPath,
	}, nil
}

// Run parses the articles in repo once and then starts observing the repo changes.
func (o *RepoObserver) Run() {
	// Walk the article file tree in repo and parse them.
	articlePath := path.Join(o.repoPath, "articles")
	filepath.WalkDir(articlePath, func(p string, d fs.DirEntry, err error) error {
		basename := d.Name()
		if err != nil ||
			d.IsDir() || // Exclude directories
			strings.HasPrefix(basename, ".") || // Exclude hidden files
			filepath.Ext(basename) != ".md" { // Exclude non-markdown files
			return nil
		}

		article, err := model.NewArticle(o.repoPath, p, d)
		if err == nil {
			o.db.Save(article)
		}

		return err
	})

	// TODO: Start observing the repo changes.
}
