package homedir

import (
	"os"
	"path/filepath"
	"strings"
)

// Expand the prefixed tilde to HOME directory.
func Expand(path string) string {
	if strings.HasPrefix(path, "~") {
		home := os.Getenv("HOME")
		rest := path[1:]
		if len(rest) == 0 || strings.HasPrefix(rest, "/") {
			return filepath.Join(home, rest)
		}
	}
	return path
}
