package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Jobs Jobs `yaml:"jobs"`

	API bool `yaml:"api"`
	UI  bool `yaml:"ui"`
}

type Job struct {
	Repo     string `yaml:"repo" json:"repo"`
	Password string `yaml:"password" json:"-"`

	Schedule string `yaml:"schedule" json:"schedule,omitempty"`
	Source   string `yaml:"source" json:"-"`
}

type Jobs map[string]Job

func Load(data []byte) (*Config, error) {
	config := Config{
		UI:  true,
		API: true,
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	for k, v := range config.Jobs {
		repo, err := filepath.Abs(v.Repo)
		if err != nil {
			return nil, err
		}
		v.Repo = repo
		config.Jobs[k] = v
	}

	return &config, nil
}

func LoadFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return Load(data)
}
