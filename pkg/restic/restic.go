package restic

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sh0rez/packup/pkg/exec"
)

type Restic struct {
	Repo     string
	Password string
}

func New(repo, pass string) (*Restic, error) {
	_, err := os.Stat(repo)
	if err != nil && !os.IsNotExist(err) {
		repo, err = filepath.Abs(repo)
		if err != nil {
			return nil, err
		}
	}

	r := &Restic{
		Repo:     repo,
		Password: pass,
	}
	return r, nil
}

const Bin = "restic"

const (
	EnvPrefix = "RESTIC_"

	EnvPass = EnvPrefix + "PASSWORD"
	EnvRepo = EnvPrefix + "REPOSITORY"
)

func (r *Restic) cmd(action string, argv []string) exec.Cmd {
	cmd := exec.Command(Bin, append([]string{action}, argv...)...)
	cmd.Env[EnvRepo] = r.Repo
	cmd.Env[EnvPass] = r.Password

	return cmd
}

type ExecError struct {
	err    error
	stderr string
}

func (e ExecError) Error() string {
	return strings.TrimSpace(fmt.Sprintf("%s.\n\n%s", e.err.Error(), e.stderr))
}

func (r *Restic) exec(action string, argv []string) ([]byte, error) {
	cmd := r.cmd(action, argv)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, ExecError{
			err:    err,
			stderr: stderr.String(),
		}
	}

	return stdout.Bytes(), nil
}

// environment is a helper type for manipulating os.Environ() more easily
type environment map[string]string

func newEnv(e []string) environment {
	env := make(environment)
	for _, s := range e {
		kv := strings.SplitN(s, "=", 2)
		env[kv[0]] = kv[1]
	}
	return env
}

func (e environment) render() []string {
	s := make([]string, 0, len(e))
	for k, v := range e {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(s)
	return s
}
