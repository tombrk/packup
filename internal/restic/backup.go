package restic

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
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

type BackupResult struct {
	Summary *BackupSummary
}

type BackupSummary struct {
	MessageType         string    `json:"message_type"`
	FilesNew            uint      `json:"files_new"`
	FilesChanged        uint      `json:"files_changed"`
	FilesUnmodified     uint      `json:"files_unmodified"`
	DirsNew             uint      `json:"dirs_new"`
	DirsChanged         uint      `json:"dirs_changed"`
	DirsUnmodified      uint      `json:"dirs_unmodified"`
	DataBlobs           int       `json:"data_blobs"`
	TreeBlobs           int       `json:"tree_blobs"`
	DataAdded           uint64    `json:"data_added"`
	DataAddedPacked     uint64    `json:"data_added_packed"`
	TotalFilesProcessed uint      `json:"total_files_processed"`
	TotalBytesProcessed uint64    `json:"total_bytes_processed"`
	TotalDuration       float64   `json:"total_duration"`
	BackupStart         time.Time `json:"backup_start"`
	BackupEnd           time.Time `json:"backup_end"`
	SnapshotID          string    `json:"snapshot_id,omitempty"`
	DryRun              bool      `json:"dry_run,omitempty"`
}

func (r *Repository) Backup(path string) (*BackupResult, error) {
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

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, newExecError(err, stderr.Bytes())
	}

	return parseBackupResult(stdout.Bytes())
}

func parseBackupResult(out []byte) (*BackupResult, error) {
	result := &BackupResult{}
	dec := json.NewDecoder(bytes.NewReader(out))

	for {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var msg Msg
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		if msg.Type != "summary" {
			continue
		}

		var summary BackupSummary
		if err := json.Unmarshal(raw, &summary); err != nil {
			return nil, err
		}
		result.Summary = &summary
	}

	return result, nil
}
