package model

import (
	"io"
	"os"
	"testing"
)

var cases = []struct {
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

func TestScanArticle(t *testing.T) {
	path := "/tmp/blogo-test-scan-article.md"
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		os.Remove(path)
	}()

	for i, c := range cases {
		// Clear file content then rewrite with text of current case.
		f.Truncate(0)
		f.Seek(0, io.SeekStart)
		f.Write([]byte(c.text))
		f.Seek(0, io.SeekStart)

		metadata, content := scanArticle(f)
		if metadata != c.metadata || content != c.content {
			t.Errorf("Failed to scan article of case[%v].", i)
		}
	}
}
