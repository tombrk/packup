package restic

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Restic struct {
	Repo     string
	Password string
}

func New(repo, pass string, job string) *Restic {
	return &Restic{
		Repo:     repo,
		Password: pass,
	}
}

func (r *Restic) cmd(action string, argv []string) *exec.Cmd {
	cmd := exec.Command("restic", append([]string{action}, argv...)...)

	env := newEnv(os.Environ())
	env["RESTIC_PASSWORD"] = r.Password
	env["RESTIC_REPOSITORY"] = r.Repo
	cmd.Env = env.render()

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
