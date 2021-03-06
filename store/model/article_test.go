package model

import (
	"os"
	"testing"
)

func TestScanArticle(t *testing.T) {
	cases := []struct {
		text     string
		metadata string
		content  string
	}{
		{
			text:     "",
			metadata: "",
			content:  "",
		},
		{
			text:     "+++\nmetadata\n+++\n",
			metadata: "metadata",
			content:  "",
		},
		{
			text:     "+++\nmetadata\n++",
			metadata: "",
			content:  "+++\nmetadata\n++",
		},
		{
			text:     "\n\ncontent\n",
			metadata: "",
			content:  "content",
		},
		{
			text:     "+++\nmetadata\n+++\n\ncontent\n",
			metadata: "metadata",
			content:  "content",
		},
		{
			text:     "\n +++\nmetadata\n\n+++\n\ncontent\n",
			metadata: "metadata",
			content:  "content",
		},
	}

	path := "/tmp/blogo-test-scan-article.md"
	for i, c := range cases {
		f, _ := os.Create(path)
		_, _ = f.Write([]byte(c.text))

		metadata, content, _ := parseArticle(path)
		if metadata != c.metadata || content != c.content {
			t.Errorf("[%v] Failed to scan article.", i)
		}

		_ = f.Close()
		_ = os.Remove(path)
	}
}
