package model

import (
	"bufio"
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/scanner"
	"time"

	"github.com/hulizhen/blogo/config"
	"github.com/hulizhen/blogo/internal/markdown"
	"github.com/hulizhen/blogo/internal/xtime"
	"github.com/hulizhen/blogo/service"

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

// NewArticle creates a Article instance with provided repo path `base`, article `path` and article `entry`.
func NewArticle(path string, entry fs.DirEntry) (article *Article, err error) {
	// Get URL slug by stripping extension of the file basename.
	basename := filepath.Base(path)
	slug := strings.TrimSuffix(basename, filepath.Ext(basename))

	// Parse article to extract metadata and content.
	metadata, content, err := parseArticle(path)
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

// parseArticle parses the *.md article file and extracts the metadata and content.
func parseArticle(path string) (metadata string, content string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
	}()

	hasMetadata := false
	removed := false
	articleScanner := bufio.NewScanner(f)
	articleScanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	for articleScanner.Scan() {
		text := strings.TrimSpace(articleScanner.Text())
		if hasMetadata && len(metadata) == 0 {
			metadata = text
		} else {
			content += text
		}
	}
	return
}

type ArticleStore struct {
	config *config.Config
	db     *sqlx.DB
}

func NewArticleStore(db *sqlx.DB, cfg *config.Config) (*ArticleStore, error) {
	s := &ArticleStore{
		config: cfg,
		db:     db,
	}

	err := s.ScanArticles()
	return s, err
}

// ScanArticles walks the artile file tree of the repo and parses them concurrently.
func (s *ArticleStore) ScanArticles() error {
	repoService := service.NewRepoService(s.config)
	if err := repoService.UpdateRepo(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	repoPath := s.config.Repository.LocalPath
	articlePath := path.Join(repoPath, "articles")
	err := filepath.WalkDir(articlePath, func(p string, d fs.DirEntry, err error) error {
		basename := d.Name()
		if err != nil ||
			d.IsDir() || // Exclude directories
			strings.HasPrefix(basename, ".") || // Exclude hidden files
			filepath.Ext(basename) != ".md" { // Exclude non-markdown files
			return nil
		}

		// Parse the articles concurrently.
		wg.Add(1)
		go func() {
			defer wg.Done()

			article, err := NewArticle(p, d)
			if err == nil {
				_, err = s.db.NamedExec(`
				REPLACE INTO articles(
					slug, title, content, preview, categories, tags, pinned, draft, published_at
				) VALUES(
					:slug, :title, :content, :preview, :categories, :tags, :pinned, :draft, :published_at
				)`,
					article,
				)
			}
			if err != nil {
				log.Panicf("Failed to parse article with error: %v.", err)
			}
		}()
		return err
	})
	wg.Wait()

	return err
}

func (s *ArticleStore) ReadArticles(limit int, offset int) ([]*Article, error) {
	rows, err := s.db.Queryx(`SELECT * FROM articles ORDER BY published_at DESC LIMIT ? OFFSET ?`, limit, offset*limit)
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
	err := s.db.Get(&count, `SELECT COUNT(*) FROM articles`)
	return count, err
}

func (s *ArticleStore) ReadArticle(slug string) (*Article, error) {
	var article Article
	if err := s.db.Get(&article, `SELECT * FROM articles WHERE slug=?`, slug); err != nil {
		return nil, err
	}
	return &article, nil
}
