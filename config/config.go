package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"reflect"

	"blogo/pkg/tilde"

	"github.com/pelletier/go-toml"
)

type website struct {
	Title        string `toml:"title"`
	Description  string `toml:"description"`
	FaviconPath  string `toml:"favicon_path" blogo:"tilde"`
	LogoPath     string `toml:"logo_path" blogo:"tilde"`
	BlogRepoPath string `toml:"blog_repo_path" blogo:"tilde"`
	TemplatePath string `toml:"template_path" blogo:"tilde"`
}

type server struct {
	Port int `toml:"port"`
}

type mysql struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// Config provides the configurations for the application.
type Config struct {
	Website website `toml:"website"`
	Server  server  `toml:"server"`
	Mysql   mysql   `toml:"mysql"`
}

var defaultConfigs = Config{
	Website: website{
		Title:        "Blogo",
		Description:  "A blog engine built with Go.",
		FaviconPath:  "~/.blogo/favicon.ico",
		LogoPath:     "~/.blogo/logo.png",
		BlogRepoPath: "~/.blogo/blog/",
		TemplatePath: "./template/",
	},
	Server: server{
		Port: 8000,
	},
	Mysql: mysql{
		Host:     "localhost",
		Port:     3306,
		Database: "blogo",
	},
}

const ConfigFilePath = "~/.blogo/config.toml"

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

// New creates a Config with the default configurations,
// which then overwritten by local custom config.toml file.
func New(p string) (*Config, error) {
	cfg := defaultConfigs

	// Parse the custom config.toml file.
	err := parseConfigFile(p, &cfg)
	if err != nil {
		return nil, err
	}

	expandTildes(&cfg)
	return &cfg, nil
}

// parseConfigFile parses the config.toml file.
func parseConfigFile(p string, cfg *Config) error {
	tilde.Expand(&p)

	f, err := os.Open(p)
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("The custom config.toml file doesn't exist, use the defaults.")
		return nil
	}
	if err == nil {
		defer f.Close()

		d := toml.NewDecoder(f)
		err = d.Decode(cfg)
	}
	return err
}
