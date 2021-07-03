package markdown

import (
	"sync"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
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
				goldmark.WithExtensions(
					highlighting.NewHighlighting(
						highlighting.WithStyle("dracula"), // theme - https://xyproto.github.io/splash/docs/all.html
						highlighting.WithFormatOptions(
							html.WithLineNumbers(true),
							html.LineNumbersInTable(true),
						),
					),
				),
			),
		}
	})
	return sharedMarkdown
}
