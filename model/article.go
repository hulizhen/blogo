package model

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pelletier/go-toml"
)

type articleMetadata struct {
	Title      string    `toml:"title"`
	Categories []string  `toml:"categories"`
	Tags       []string  `toml:"tags"`
	Top        bool      `toml:"top"`
	Draft      bool      `toml:"drfat"`
	Date       time.Time `toml:"date"`
}

type Article struct {
	ID          int64     `gorm:"primarykey"`
	Slug        string    `gorm:"slug"`
	Title       string    `gorm:"title"`
	Content     string    `gorm:"content"`
	Preview     string    `gorm:"preview"`
	Categories  string    `gorm:"categories"`
	Tags        string    `gorm:"tags"`
	Top         bool      `gorm:"top"`
	Draft       bool      `gorm:"draft"`
	PublishedTS time.Time `gorm:"published_ts"`
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
	m, c := scanArticle(path)
	if err != nil {
		return
	}

	// Extract preview from content.
	preview := ""
	count := 3
	strs := strings.SplitN(c, previewDelimiter, count)
	if len(strs) == count {
		preview = strings.TrimSpace(strs[1])
		c = strs[0] + strs[1] + strs[2]
	}

	// Get URL slug by stripping extension of the file basename.
	basename := filepath.Base(path)
	slug := strings.TrimSuffix(basename, filepath.Ext(basename))

	// Parse metadata and fill article.
	article = &Article{
		ID:      id,
		Slug:    slug,
		Content: c,
		Preview: preview,
	}
	am := &articleMetadata{}
	err = toml.Unmarshal([]byte(m), am)
	if err != nil {
		return
	}
	article.updateMetadata(am)

	return
}

// updateMetadata updates the article model with the extracted metadata.
func (a *Article) updateMetadata(am *articleMetadata) {
	a.Title = am.Title
	a.Categories = strings.Join(am.Categories, categoryDelimiter)
	a.Tags = strings.Join(am.Tags, tagDelimiter)
	a.Top = am.Top
	a.Draft = am.Draft
	a.PublishedTS = am.Date
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

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
