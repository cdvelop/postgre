package postgre

import (
	"os"
	"time"
)

func findOldestFile(dir string) (oldestFile os.FileInfo, total_files int, err error) {
	dir_files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	total_files = len(dir_files)

	oldestTime := time.Now()
	for _, dir_file := range dir_files {

		file, er := dir_file.Info()
		if er != nil {
			err = er
			return
		}

		if file.Mode().IsRegular() && file.ModTime().Before(oldestTime) {
			oldestFile = file
			oldestTime = file.ModTime()
		}
	}

	if oldestFile == nil {
		err = os.ErrNotExist
	}
	return
}
