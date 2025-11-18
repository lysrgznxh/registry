package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func Clean(src string, timeoutSeconds int64) error {
	now := time.Now()
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && int64(now.Sub(info.ModTime()).Seconds()) > timeoutSeconds {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
