package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Jobs Jobs `yaml:"jobs"`
}

type Job struct {
	Repo     string `yaml:"repo" json:"repo"`
	Password string `yaml:"password" json:"-"`
	Schedule string `yaml:"schedule" json:"schedule,omitempty"`
}

type Jobs map[string]Job

func Load(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
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
