package restic

import (
	"io"
	"os"
)

func (r Restic) Dump(w io.Writer, snapshot string, filepath string) error {
	cmd := r.cmd("dump", []string{snapshot, filepath})
	cmd.Stderr = os.Stderr
	cmd.Stdout = w

	return cmd.Run()
}
