package config

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

type Source struct {
	Path string `yaml:"path,omitempty"`
	Exec string `yaml:"exec,omitempty"`
}

func (s Source) Dir() (string, error) {
	if s.Empty() {
		return "", ErrorNoSource
	}

	path := s.Path
	if s.Exec != "" {
		var err error
		path, err = s.exec()
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}

func (s Source) Empty() bool {
	return s.Exec == "" && s.Path == ""
}

var ErrorNoSource = errors.New("No source available, as 'path' and 'exec' both unset")

func (s Source) exec() (string, error) {
	cmd := exec.Command("/bin/sh", "-c", s.Exec)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
