package restic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	MessageType string `json:"message_type,omitempty"`
}

func (r *Repository) Files(snapshot string, path string, recursive bool) ([]File, error) {
	args := []string{snapshot, path, "--json", "--no-lock"}
	if recursive {
		args = append(args, "--recursive")
	}

	out, err := r.exec("ls", args)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(out))

	var files []File
	for {
		var f File
		if err := dec.Decode(&f); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// restic ls --json emits a leading snapshot object and, in newer
		// versions, explicit message_type values. Keep only actual file nodes.
		if f.MessageType != "" && f.MessageType != "node" {
			continue
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
