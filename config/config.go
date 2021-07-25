package config

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"os"
	"reflect"

	"github.com/hulizhen/blogo/pkg/tilde"

	"github.com/pelletier/go-toml"
)

type website struct {
	Title           string `toml:"title"`
	Description     string `toml:"description"`
	Author          string `toml:"author"`
	SinceYear       int    `toml:"since_year"`
	ArticlePageSize int    `toml:"article_page_size"`
	FaviconPath     string `toml:"favicon_path" blogo:"tilde"`
	LogoPath        string `toml:"logo_path" blogo:"tilde"`
	TemplatePath    string `toml:"template_path" blogo:"tilde"`
}

type repository struct {
	LocalPath string `toml:"local_path" blogo:"tilde"`
	RemoteURL string `toml:"remote_url"`
	Branch    string `toml:"branch"`
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
	Website    website    `toml:"website"`
	Repository repository `toml:"repository"`
	Server     server     `toml:"server"`
	Mysql      mysql      `toml:"mysql"`
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

//go:embed config.toml
var defaultConfigFile embed.FS

// New creates a Config with the default config.toml file,
// which then overwritten by local custom config.toml file.
func New() (cfg *Config, err error) {
	cfg = new(Config)

	// Parse the default config.toml file.
	f, err := defaultConfigFile.Open("config.toml")
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	d := toml.NewDecoder(f)
	err = d.Decode(cfg)
	if err != nil {
		return
	}

	// Parse the custom config.toml file.
	var customConfigPath = "~/.blogo/config.toml"
	tilde.Expand(&customConfigPath)
	f, err = os.Open(customConfigPath)
	if errors.Is(err, fs.ErrNotExist) {
		log.Printf("The custom config.toml file doesn't exist, use the defaults.\n")
		err = nil
	} else {
		defer func() {
			_ = f.Close()
		}()
		d = toml.NewDecoder(f)
		err = d.Decode(cfg)
	}
	if err != nil {
		return
	}

	expandTildes(cfg)
	return
}
