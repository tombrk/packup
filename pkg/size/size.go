package size

import (
	"os"
	"path/filepath"
)

// Size reports the size of that directory
func Size(dir string) (int64, error) {
	n := int64(0)
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			n += fi.Size()
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return n, nil
}
