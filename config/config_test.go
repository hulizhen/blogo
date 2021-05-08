package config

import (
	"testing"
)

var cases = []struct {
	path string
	port int
}{
	{"~/test/config.toml", 8000},
	{"./testdata/config.toml", 8080},
}

func TestNew(t *testing.T) {
	for _, c := range cases {
		cfg := new(c.path)
		if cfg.Server.Port != c.port {
			t.Errorf("the port in config should be '%d'", c.port)
		}
	}
}
