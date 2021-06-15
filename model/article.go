package model

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pelletier/go-toml"
)

type articleMetadata struct {
	Title string
	Tags  []string
	Draft bool
	Date  time.Time
}

type Article struct {
	ID          int64     `gorm:"primarykey"`
	Title       string    `gorm:"title"`
	Content     string    `gorm:"content"`
	Categories  string    `gorm:"categories"`
	Tags        string    `gorm:"tags"`
	Draft       bool      `gorm:"draft"`
	PublishedTS time.Time `gorm:"published_ts"`
}

const metadataDelimiter = "+++"

// New creates a Article instance with provided repo path `base`, article `path` and article `entry`.
func NewArticle(base string, path string, entry fs.DirEntry) (article *Article, err error) {
	// Get categories.
	parent := filepath.Dir(path)
	categories := strings.TrimLeft(strings.TrimPrefix(parent, base), "/")

	// Generate ID with birth timestamp.
	id := int64(0)
	info, err := entry.Info()
	if err == nil {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			err = fmt.Errorf("failed to get stat information for article path: %v", path)
			return
		}
		id = stat.Birthtimespec.Nano()
	}

	// Scan article to extract metadata and content.
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	metadata, content := scanArticle(f)
	if err != nil {
		return
	}

	// Parse metadata and fill article.
	article = &Article{
		ID:         id,
		Content:    content,
		Categories: categories,
	}
	am := &articleMetadata{}
	err = toml.Unmarshal([]byte(metadata), am)
	if err != nil {
		return
	}
	article.fillMetadata(am)

	return
}

func (a *Article) fillMetadata(metadata *articleMetadata) {
	a.Title = metadata.Title
	a.Tags = strings.Join(metadata.Tags, ",")
	a.Draft = metadata.Draft
	a.PublishedTS = metadata.Date
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func scanArticle(f *os.File) (metadata string, content string) {
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
		strs := strings.SplitN(string(data), metadataDelimiter, 3)
		if len(strs) == 3 && len(strs[0]) == 0 {
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
