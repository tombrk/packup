package restic

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sh0rez/packup/pkg/exec"
)

type Repository struct {
	Addr     string
	Password string
}

func Open(addr, pass string) (*Repository, error) {
	_, err := os.Stat(addr)
	if err != nil && !os.IsNotExist(err) {
		addr, err = filepath.Abs(addr)
		if err != nil {
			return nil, err
		}
	}

	r := &Repository{
		Addr:     addr,
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

func (r *Repository) cmd(action string, argv []string) exec.Cmd {
	cmd := exec.Command(Bin, append([]string{action}, argv...)...)
	cmd.Env[EnvRepo] = r.Addr
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

func (r *Repository) exec(action string, argv []string) ([]byte, error) {
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
