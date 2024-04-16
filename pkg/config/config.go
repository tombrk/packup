package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Jobs Jobs `yaml:"jobs"`
}

type Job struct {
	Repo     string `yaml:"repo" json:"repo"`
	Password string `yaml:"password" json:"-"`
}

type Jobs map[string]Job

func Parse(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var errs error
	for name, job := range config.Jobs {
		repo, err := filepath.Abs(job.Repo)
		errs = errors.Join(errs, err)
		job.Repo = repo
		config.Jobs[name] = job
	}
	return &config, errs
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}
