package config

import (
	"testing"
)

var cases = []struct {
	paths []string
	port  int
}{
	{[]string{"/tmp/config.toml", "~/test/config.toml"}, 8000},
	{[]string{"/tmp/config.toml", "testdata/config.toml"}, 8080},
}

func TestNew(t *testing.T) {
	for _, c := range cases {
		cfg := new(c.paths, defaultConfigs)
		if cfg.Server.Port != c.port {
			t.Errorf("the port in config should be '%d'", c.port)
		}
	}
}
