package restic

import (
	"encoding/json"
	"fmt"
	"time"
)

var ErrNoSnapshots = fmt.Errorf("Repository has no snapshots. Please run a backup first")

type Snapshot struct {
	ID       string           `json:"id"`
	Hostname string           `json:"hostname"`
	Time     time.Time        `json:"time"`
	Summary  *SnapshotSummary `json:"summary,omitempty"`
}

type SnapshotSummary struct {
	BackupStart         time.Time `json:"backup_start"`
	BackupEnd           time.Time `json:"backup_end"`
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
}

func (r *Repository) Snapshots() ([]Snapshot, error) {
	out, err := r.exec("snapshots", []string{"--json", "--no-lock", "--no-cache"})
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	if err := json.Unmarshal(out, &snapshots); err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, ErrNoSnapshots
	}

	return snapshots, nil
}
