package pagination

import (
	"fmt"
	"html"
	"math"
)

type Page struct {
	Label    string
	Text     string
	Href     string
	Disabled bool
}

type Pagination struct {
	Total   int
	Current int
	Path    string
}

const DefaultPageSize = 10

func New(count, size, current int, path string) *Pagination {
	total := int(math.Ceil(float64(count) / float64(size)))
	return &Pagination{
		Total:   total,
		Current: current,
		Path:    path,
	}
}

func (p *Pagination) Previous() Page {
	offset := p.Current - 1
	if offset > p.Total {
		offset = p.Total
	}
	return Page{
		Label:    "Previous",
		Text:     html.UnescapeString("&lsaquo;"),
		Href:     p.buildHref(offset),
		Disabled: p.Current <= 1,
	}
}

func (p *Pagination) Next() Page {
	offset := p.Current + 1
	if offset < 1 {
		offset = 1
	}
	return Page{
		Label:    "Next",
		Text:     html.UnescapeString("&rsaquo;"),
		Href:     p.buildHref(offset),
		Disabled: p.Current >= p.Total,
	}
}

func (p *Pagination) First() Page {
	return Page{
		Label:    "First",
		Text:     html.UnescapeString("&laquo;"),
		Href:     p.buildHref(1),
		Disabled: p.Current == 1,
	}
}

func (p *Pagination) Last() Page {
	return Page{
		Label:    "Last",
		Text:     html.UnescapeString("&raquo;"),
		Href:     p.buildHref(p.Total),
		Disabled: p.Current == p.Total,
	}
}

func (p *Pagination) buildHref(offset int) string {
	return fmt.Sprintf("%v?offset=%d", p.Path, offset)
}
