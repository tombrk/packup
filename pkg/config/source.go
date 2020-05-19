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
	if s.Exec != "" {
		return s.exec()
	}

	if s.Path != "" {
		return s.Path, nil
	}

	return "", ErrorNoSource
}

var ErrorNoSource = errors.New("No source available, as 'path' and 'exec' both unset")

func (s Source) exec() (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s.Exec)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
