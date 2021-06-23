package observer

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"blogo/model"

	"github.com/jmoiron/sqlx"
)

type RepoObserver struct {
	db       *sqlx.DB
	repoPath string
}

func NewRepoObserver(db *sqlx.DB, repoPath string) (*RepoObserver, error) {
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
			_, err = o.db.NamedExec(`
			INSERT INTO article(
				id, slug, title, content, preview, categories, tags, top, draft, published_ts
			) VALUES(
				:id, :slug, :title, :content, :preview, :categories, :tags, :top, :draft, :published_ts
			) ON DUPLICATE KEY UPDATE
				id = :id, slug = :slug, title = :title, content = :content, preview = :preview, categories = :categories, tags = :tags, top = :top, draft = :draft, published_ts = :published_ts
			`, article)
		}
		return err
	})

	// TODO: Start observing the repo changes.
}
