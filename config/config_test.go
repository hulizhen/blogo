package config

import (
	"os"
	"strings"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	cases := []struct {
		path string
		port int
	}{
		{"~/test/config.toml", 8000},
		{"./testdata/config.toml", 8080},
	}

	for _, c := range cases {
		cfg := Config{Server: server{Port: 8000}}
		_ = parseConfigFile(c.path, &cfg)
		if cfg.Server.Port != c.port {
			t.Errorf("The port in config should be '%d'", c.port)
		}
	}
}

func TestExpandTildes(t *testing.T) {
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
	for _, c := range cases {
		expandTildes(&c)
		if strings.HasPrefix(c.Embeded.InnerString, home) || strings.HasPrefix(c.OuterString, home) {
			t.Errorf("The tildes in the field without tag should not be expanded.")
		}
		if !strings.HasPrefix(c.Embeded.InnerPath, home) || !strings.HasPrefix(c.OuterPath, home) {
			t.Errorf("The tildes in the field with tag should be expanded.")
		}
	}
}
