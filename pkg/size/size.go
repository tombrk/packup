package size

import (
	"bytes"
	"strconv"
	"strings"

	"os/exec"
)

// Size reports the size of that directory
func Size(dir string) (int64, error) {
	cmd := exec.Command("du", "--bytes", "--summarize", dir)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	s := strings.Fields(out.String())[0]
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}
