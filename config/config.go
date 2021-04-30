package config

import (
	"fmt"
	"os"

	"github.com/hulizhen/blogo/pkg/tilde"
	"github.com/pelletier/go-toml"
)

type website struct {
	FaviconPath  string `toml:"favicon_path"`
	BlogRepoPath string `toml:"blog_repo_path"`
	Title        string `toml:"title"`
	Description  string `toml:"description"`
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
		FaviconPath:  "~/.blogo/favicon.ico",
		BlogRepoPath: "~/.blogo/blog",
		Title:        "Blogo",
		Description:  "A blog engine built with Go.",
	},
	Server: server{
		Port: 8000,
	},
}

var configFilePaths []string = []string{
	"config.toml",
	"~/.blogo/config.toml",
}

var sharedConfig *Config

// SharedConfig always returns a singleton of the Config instance
// to share in the whole application.
func SharedConfig() *Config {
	if sharedConfig == nil {
		sharedConfig = new(configFilePaths, defaultConfigs)
	}
	return sharedConfig
}

// new creates a Config with the provided default configurations,
// which then overwritten by local custom config file.
func new(paths []string, defaults Config) *Config {
	cfg := defaults
	f, err := openConfigFile(paths)
	if f != nil {
		d := toml.NewDecoder(f)
		err = d.Decode(&cfg)
	}
	if err != nil {
		fmt.Printf("Failed to open config file with error: %v, use the defaults.\n", err.Error())
	}
	tilde.Expand(&cfg.Website.FaviconPath)
	tilde.Expand(&cfg.Website.BlogRepoPath)
	return &cfg
}

func openConfigFile(paths []string) (*os.File, error) {
	for _, p := range paths {
		tilde.Expand(&p)
		if file, err := os.Open(p); err == nil {
			return file, nil
		}
	}
	return nil, fmt.Errorf("the config.toml file is not found in the search paths '%v'", paths)
}
