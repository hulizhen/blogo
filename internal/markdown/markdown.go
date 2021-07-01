package markdown

import (
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

type Markdown struct {
	goldmark.Markdown
}

var once sync.Once
var sharedMarkdown *Markdown

// SharedMarkdown always returns a singleton of the Markdown instance
// to share in the whole application.
func SharedMarkdown() *Markdown {
	once.Do(func() {
		sharedMarkdown = &Markdown{
			goldmark.New(
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
				),
			),
		}
	})
	return sharedMarkdown
}
