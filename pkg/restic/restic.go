package restic

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Restic struct {
}

func New() *Restic {
	return &Restic{}
}

type Snapshot struct {
	ID       string    `json:"id"`
	Hostname string    `json:"hostname"`
	Time     time.Time `json:"time"`
}

func (r *Restic) Snapshots() ([]Snapshot, error) {
	out, err := r.exec("snapshots", []string{"--json"})
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	if err := json.Unmarshal(out, &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

func (r *Restic) Files(snapshot string, path string, recursive bool) ([]File, error) {
	args := []string{snapshot, path, "--json"}
	if recursive {
		args = append(args, "--recursive")
	}

	out, err := r.exec("ls", args)
	if err != nil {
		return nil, err
	}

	lines := strings.TrimSuffix(string(out), "\n")

	var files []File
	for _, l := range strings.Split(lines, "\n") {
		var f File
		if err := json.Unmarshal([]byte(l), &f); err != nil {
			return nil, err
		}

		if !(f.Type == "file" || f.Type == "dir") {
			continue
		}

		files = append(files, f)
	}

	return files, nil
}

func (r *Restic) exec(action string, argv []string) ([]byte, error) {
	cmd := exec.Command("restic", append([]string{action}, argv...)...)
	cmd.Stderr = os.Stderr

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
