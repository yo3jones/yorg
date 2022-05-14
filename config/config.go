package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataRoot string `yaml:"dataRoot"`
}

type userHomeDirer struct{}

type UserHomeDirer interface {
	UserHomeDir() (string, error)
}

func (_ *userHomeDirer) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

type accessor struct {
	userHomeDirer UserHomeDirer
}

type Accessor interface {
	Access() (*Config, error)
}

func (a *accessor) Access() (*Config, error) {
	home, err1 := a.userHomeDirer.UserHomeDir()
	if err1 != nil {
		return nil, err1
	}

	path := filepath.Join(home, ".yorg.yaml")
	f, err2 := os.Open(path)
	if err2 != nil {
		return nil, err2
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)

	var config Config
	err3 := decoder.Decode(&config)
	if err3 != nil {
		return nil, err3
	}

	return &config, nil
}

func Load() (*Config, error) {
	a := &accessor{
		userHomeDirer: &userHomeDirer{},
	}
	return a.Access()
}
