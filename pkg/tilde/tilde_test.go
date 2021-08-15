package tilde

import (
	"os"
	"path/filepath"
	"testing"
)

var home = os.Getenv("HOME")

func TestExpandTilde(t *testing.T) {
	cases := []struct {
		path   string
		result string
	}{
		{"~~", "~~"},
		{"~", home},
		{"~/", filepath.Join(home, "/")},
		{"~/TestPath", filepath.Join(home, "TestPath")},
		{"TestA/~", "TestA/~"},
		{"TestA/TestB", "TestA/TestB"},
		{"/TestA/TestB", "/TestA/TestB"},
	}

	for i, c := range cases {
		Expand(&c.path)
		if c.path != c.result {
			t.Errorf("[%v] Expand '%v' should result in '%v'.", i, c.path, c.result)
		}
	}
}
