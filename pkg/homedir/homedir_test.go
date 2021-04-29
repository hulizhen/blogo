package homedir

import (
	"os"
	"path/filepath"
	"testing"
)

var home = os.Getenv("HOME")

var cases = []struct {
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

func TestExpand(t *testing.T) {
	for _, c := range cases {
		if Expand(c.path) != c.result {
			t.Errorf("expand '%v' should result in '%v'", c.path, c.result)
		}
	}
}
