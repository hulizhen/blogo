package model

import (
	"blogo/config"
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type AboutStore struct {
	About map[string]interface{}
}

func NewAboutStore(cfg *config.Config) (*AboutStore, error) {
	p := filepath.Join(cfg.Website.BlogRepoPath, "about.toml")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var about map[string]interface{}
	if err = toml.Unmarshal(b, &about); err != nil {
		return nil, err
	}

	return &AboutStore{About: about}, nil
}
