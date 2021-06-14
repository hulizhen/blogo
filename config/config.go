package config

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/hulizhen/blogo/pkg/tilde"
	"github.com/pelletier/go-toml"
)

type website struct {
	Name         string `toml:"name"`
	Description  string `toml:"description"`
	FaviconPath  string `toml:"favicon_path" blogo:"tilde"`
	LogoPath     string `toml:"logo_path" blogo:"tilde"`
	BlogRepoPath string `toml:"blog_repo_path" blogo:"tilde"`
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

// TODO: Pass the config instance around instead of creating a global share instance.
// SharedConfig always returns a singleton of the Config instance
// to share in the whole application.
func SharedConfig() *Config {
	once.Do(func() {
		sharedConfig = new(customConfigFilePath)
	})
	return sharedConfig
}

// expandTildes expands tildes of the path strings in the Config instance recursively.
func expandTildes(x interface{}) {
	v := reflect.Indirect(reflect.ValueOf(x))
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := v.Type().Field(i)
		value, ok := tf.Tag.Lookup("blogo")
		if ok && value == "tilde" {
			ptr := vf.Addr().Interface().(*string)
			tilde.Expand(ptr)
		}
		if tf.Type.Kind() == reflect.Struct {
			expandTildes(vf.Addr().Interface())
		}
	}
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

	expandTildes(&cfg)
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
