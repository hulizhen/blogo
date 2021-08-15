package config

import (
	"os"
	"strings"
	"testing"
)

func TestExpandTildesx(t *testing.T) {
	type embeded struct {
		InnerString string
		InnerPath   string `blogo:"tilde"`
	}
	cases := []struct {
		Embeded     embeded
		OuterString string
		OuterPath   string `blogo:"tilde"`
	}{
		{
			Embeded: embeded{
				InnerString: "~/test/inner/string",
				InnerPath:   "~/test/inner/path",
			},
			OuterString: "~/test/outer/string",
			OuterPath:   "~/test/outer/path",
		},
	}
	home := os.Getenv("HOME")
	for i, c := range cases {
		expandTildes(&c)
		if strings.HasPrefix(c.Embeded.InnerString, home) || strings.HasPrefix(c.OuterString, home) {
			t.Errorf("[%v] The tildes in the field without tag should not be expanded.", i)
		}
		if !strings.HasPrefix(c.Embeded.InnerPath, home) || !strings.HasPrefix(c.OuterPath, home) {
			t.Errorf("[%v] The tildes in the field with tag should be expanded.", i)
		}
	}
}
