package model

import (
	"blogo/config"
	"blogo/internal/markdown"
	"blogo/internal/xtime"
	"bufio"
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/scanner"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pelletier/go-toml"
)

type articleMetadata struct {
	Title      string    `toml:"title"`
	Categories []string  `toml:"categories"`
	Tags       []string  `toml:"tags"`
	Pinned     bool      `toml:"pinned"`
	Draft      bool      `toml:"drfat"`
	Date       time.Time `toml:"date"`
}

type Article struct {
	ID          int64         `db:"id"`
	Slug        string        `db:"slug"`
	Title       string        `db:"title"`
	Content     template.HTML `db:"content"`
	Preview     template.HTML `db:"preview"`
	Categories  string        `db:"categories"`
	Tags        string        `db:"tags"`
	Pinned      bool          `db:"pinned"`
	Draft       bool          `db:"draft"`
	PublishedAt time.Time     `db:"published_at"`
}

const (
	metadataDelimiter = "+++\n"
	previewDelimiter  = "^^^\n"
	categoryDelimiter = "/"
	tagDelimiter      = ","
)

// New creates a Article instance with provided repo path `base`, article `path` and article `entry`.
func NewArticle(base string, path string, entry fs.DirEntry) (article *Article, err error) {
	// Generate ID with birth timestamp.
	id := int64(0)
	info, err := entry.Info()
	if err == nil {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			log.Printf("failed to get stat information for article path: %v", path)
			return
		}
		id = stat.Birthtimespec.Nano()
	}

	// Scan article to extract metadata and content.
	metadata, content := scanArticle(path)
	if err != nil {
		return
	}

	// Extract preview from content.
	preview := ""
	count := 3
	strs := strings.SplitN(content, previewDelimiter, count)
	if len(strs) == count {
		preview = strings.TrimSpace(strs[1])
		content = strs[0] + strs[1] + strs[2]
	}

	// Get URL slug by stripping extension of the file basename.
	basename := filepath.Base(path)
	slug := strings.TrimSuffix(basename, filepath.Ext(basename))

	// Parse the content and preview with markdown parser.
	var buf bytes.Buffer
	if err = markdown.SharedMarkdown().Convert([]byte(content), &buf); err != nil {
		return
	}
	content = buf.String()
	buf.Reset()
	if err = markdown.SharedMarkdown().Convert([]byte(preview), &buf); err != nil {
		return
	}
	preview = buf.String()

	// Parse metadata and fill article.
	article = &Article{
		ID:      id,
		Slug:    slug,
		Content: template.HTML(content),
		Preview: template.HTML(preview),
	}
	am := &articleMetadata{}
	err = toml.Unmarshal([]byte(metadata), am)
	if err != nil {
		return
	}
	article.updateMetadata(am)

	return
}

func (a *Article) ShortPublicationDate() string {
	return xtime.ShortFormat(a.PublishedAt)
}

func (a *Article) Href() string {
	return filepath.Join("/articles", a.Slug)
}

// updateMetadata updates the article model with the extracted metadata.
func (a *Article) updateMetadata(am *articleMetadata) {
	a.Title = am.Title
	a.Categories = strings.Join(am.Categories, categoryDelimiter)
	a.Tags = strings.Join(am.Tags, tagDelimiter)
	a.Pinned = am.Pinned
	a.Draft = am.Draft
	a.PublishedAt = am.Date
}

func isWhitespace(c byte) bool {
	return scanner.GoWhitespace&(1<<c) != 0
}

// scanArticle scans the *.md article file and extracts the metadata and content.
func scanArticle(path string) (metadata string, content string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	hasMetadata := false
	removed := false
	scanner := bufio.NewScanner(f)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return
		}
		if len(metadata) > 0 {
			advance = len(data)
			token = data
			return
		}

		// Remove the whitespaces at the file beginning if they exist.
		if !removed {
			removed = true
			i := 0
			for i < len(data) && isWhitespace(data[i]) {
				i++
			}
			if i > 0 {
				return i, nil, nil
			}
		}

		// Extract metadata with delimiter.
		// The metadata surrounded by delimiter should be at the beginning of article, with some whitespaces, if any.
		count := 3
		strs := strings.SplitN(string(data), metadataDelimiter, count)
		if len(strs) == count && len(strs[0]) == 0 {
			hasMetadata = true
			return len(strs[1]) + 2*len(metadataDelimiter), []byte(strs[1]), nil
		} else {
			if atEOF {
				return len(data), data, nil
			} else {
				return
			}
		}
	})
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if hasMetadata && len(metadata) == 0 {
			metadata = text
		} else {
			content += text
		}
	}
	return
}

type ArticleStore struct {
	db *sqlx.DB
}

func NewArticleStore(db *sqlx.DB, cfg *config.Config) (*ArticleStore, error) {
	repoPath := cfg.Website.BlogRepoPath
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, err
	}

	// Walk the article file tree in repo and parse them.
	articlePath := path.Join(repoPath, "articles")
	err := filepath.WalkDir(articlePath, func(p string, d fs.DirEntry, err error) error {
		basename := d.Name()
		if err != nil ||
			d.IsDir() || // Exclude directories
			strings.HasPrefix(basename, ".") || // Exclude hidden files
			filepath.Ext(basename) != ".md" { // Exclude non-markdown files
			return nil
		}

		article, err := NewArticle(repoPath, p, d)
		if err == nil {
			_, err = db.NamedExec(`
				REPLACE INTO article(
					id, slug, title, content, preview, categories, tags, pinned, draft, published_at
				) VALUES(
					:id, :slug, :title, :content, :preview, :categories, :tags, :pinned, :draft, :published_at
				)`,
				article,
			)
		}
		return err
	})

	// TODO: Start listening repo webhook to rescan the articles.

	return &ArticleStore{db: db}, err
}

func (s *ArticleStore) ReadArticles(limit int, offset int) ([]*Article, error) {
	rows, err := s.db.Queryx(`SELECT * FROM article ORDER BY published_at DESC LIMIT ? OFFSET ?`, limit, offset*limit)
	if err != nil {
		return nil, err
	}

	var articles []*Article
	for rows.Next() {
		var article Article
		err = rows.StructScan(&article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}
	return articles, nil
}

func (s *ArticleStore) ReadArticleCount() (int, error) {
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM article`)
	return count, err
}

func (s *ArticleStore) ReadArticle(slug string) (*Article, error) {
	var article Article
	if err := s.db.Get(&article, `SELECT * FROM article WHERE slug=?`, slug); err != nil {
		return nil, err
	}
	return &article, nil
}
