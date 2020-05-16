package restic

import (
	"encoding/json"
	"fmt"
	"time"
)

var ErrNoSnapshots = fmt.Errorf("Repository has no snapshots. Please run a backup first")

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

	if len(snapshots) == 0 {
		return nil, ErrNoSnapshots
	}

	return snapshots, nil
}
