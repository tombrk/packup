package restic

import (
	"bytes"
	"os"
	"os/exec"
)

type Restic struct {
	Repo     string
	Password string
}

func New(repo, pass string) *Restic {
	return &Restic{
		Repo:     repo,
		Password: pass,
	}
}

func (r *Restic) cmd(action string, argv []string) *exec.Cmd {
	return exec.Command("restic", append([]string{action}, argv...)...)
}

func (r *Restic) exec(action string, argv []string) ([]byte, error) {
	cmd := r.cmd(action, argv)
	cmd.Stderr = os.Stderr

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
