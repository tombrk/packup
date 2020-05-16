package restic

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
		if l == "" {
			continue
		}

		var f File
		if err := json.Unmarshal([]byte(l), &f); err != nil {
			return nil, err
		}

		if !(f.Type == "file" || f.Type == "dir") {
			continue
		}

		if f.Path == path && f.Type != "file" {
			continue
		}

		files = append(files, f)
	}

	if len(files) == 0 {
		return nil, ErrNoFiles{Snapshot: snapshot, Path: path}
	}

	return files, nil
}

type ErrNoFiles struct {
	Snapshot string
	Path     string
}

func (e ErrNoFiles) Error() string {
	return fmt.Sprintf(`Listing path '%s' of snapshot '%s' did not return any files.
This usually means the snapshot does not exist`, e.Path, e.Snapshot)
}
