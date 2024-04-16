package restic

import (
	"io"
	"os"
	"strings"
)

func (r Repository) Dump(w io.Writer, snapshot string, filepath string) error {
	path := strings.TrimPrefix(filepath, "/")

	cmd := r.cmd("dump", []string{snapshot, path})
	cmd.Stderr = os.Stderr
	cmd.Stdout = w

	return cmd.Run()
}
