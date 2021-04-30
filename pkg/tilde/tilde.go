package tilde

import (
	"os"
	"path/filepath"
	"strings"
)

// Expand expands the prefixed tilde to HOME directory.
func Expand(path *string) {
	if strings.HasPrefix(*path, "~") {
		home := os.Getenv("HOME")
		rest := (*path)[1:]
		if len(rest) == 0 || strings.HasPrefix(rest, "/") {
			*path = filepath.Join(home, rest)
		}
	}
}
