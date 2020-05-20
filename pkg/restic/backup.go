package restic

import (
	"bytes"
)

type BackupStatus struct {
	FilesDiscovered int `json:"total_files"`
	FilesWritten    int `json:"files_done"`

	BytesDiscovered int `json:"total_bytes"`
	BytesWritten    int `json:"bytes_done"`
}

type BackupSummary struct {
	BytesAdded int `json:"data_added"`
}

func (r *Restic) Backup(path string, status func(BackupStatus)) error {
	cmd := r.cmd("backup", []string{"--json", "."})

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
