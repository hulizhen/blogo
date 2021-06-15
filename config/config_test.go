package config

import (
	"os"
	"strings"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	cases := []struct {
		text string
		port int
	}{
		{"", 8000},
		{"[server]\nport = 8080\n", 8080},
	}

	path := "/tmp/blogo-test-config.toml"
	for i, c := range cases {
		f, _ := os.Create(path)
		f.Write([]byte(c.text))
		f.Close()

		cfg := Config{Server: server{Port: 8000}}
		if err := parseConfigFile(path, &cfg); err != nil {
			t.Errorf("[%v] Failed to parse config file with error: %v", i, err)
		}
		if cfg.Server.Port != c.port {
			t.Errorf("[%v] The port in config should be '%d'.", i, c.port)
		}
	}
	os.Remove(path)
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
