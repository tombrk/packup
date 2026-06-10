package restic

import (
	"bytes"
	"io"
	"strings"
)

func (r Repository) Dump(w io.Writer, snapshot string, filepath string) error {
	path := strings.TrimPrefix(filepath, "/")

	cmd := r.cmd("dump", []string{snapshot, path, "--no-lock"})
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		return newExecError(err, stderr.Bytes())
	}
	return nil
}
