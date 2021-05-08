package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/hulizhen/blogo/pkg/tilde"
	"github.com/pelletier/go-toml"
)

type website struct {
	Name         string `toml:"name"`
	Description  string `toml:"description"`
	FaviconPath  string `toml:"favicon_path"`
	LogoPath     string `toml:"logo_path"`
	BlogRepoPath string `toml:"blog_repo_path"`
}

type server struct {
	Port int `toml:"port"`
}

// Config provides the configurations for the application.
type Config struct {
	Website website `toml:"website"`
	Server  server  `toml:"server"`
}

var defaultConfigs = Config{
	Website: website{
		Name:         "Blogo",
		Description:  "A blog engine built with Go.",
		FaviconPath:  "~/.blogo/favicon.ico",
		LogoPath:     "~/.blogo/logo.png",
		BlogRepoPath: "~/.blogo/blog",
	},
	Server: server{
		Port: 8000,
	},
}

var customConfigFilePath = "~/.blogo/config.toml"

var once sync.Once
var sharedConfig *Config

// SharedConfig always returns a singleton of the Config instance
// to share in the whole application.
func SharedConfig() *Config {
	once.Do(func() {
		sharedConfig = new(customConfigFilePath)
	})
	return sharedConfig
}

// new creates a Config with the default configurations,
// which then overwritten by local custom config.toml file.
func new(p string) *Config {
	cfg := defaultConfigs

	// Parse the custom config.toml file.
	err := parseConfigFile(p, &cfg)
	if err != nil {
		fmt.Printf("Failed to parse the custom configurations with error: %v, use the defaults.\n", err)
	}

	// Expand the tilde in path strings.
	tilde.Expand(&cfg.Website.FaviconPath)
	tilde.Expand(&cfg.Website.LogoPath)
	tilde.Expand(&cfg.Website.BlogRepoPath)

	return &cfg
}

// parseConfigFile parses the config.toml file.
func parseConfigFile(p string, cfg *Config) error {
	tilde.Expand(&p)
	f, err := os.Open(p)
	if f != nil && err == nil {
		d := toml.NewDecoder(f)
		err = d.Decode(cfg)
	}
	return err
}
