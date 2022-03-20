package restic

import (
	"bytes"
)

type Msg struct {
	Type string `json:"message_type"`
}

type BackupStatus struct {
	FilesTotal int `json:"total_files"`
	FilesDone  int `json:"files_done"`

	BytesTotal int `json:"total_bytes"`
	BytesDone  int `json:"bytes_done"`
}

func (r *Restic) Backup(path string) error {
	if path == "" {
		panic("path must not be empty")
	}

	cmd := r.cmd("backup", []string{"--json", "."})

	// special case for backup:
	// we don't want restic to only backup file contents, not the directory those were obtained from.
	// thus, we cd to the dir and explicitely set the source to "."
	cmd.Dir = path

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ExecError{
			err:    err,
			stderr: stderr.String(),
		}
	}

	return nil
}
