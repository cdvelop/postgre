package postgre

import (
	"os"
	"time"
)

func findOldestFile(dir string) (oldestFile os.FileInfo, total_files int, err string) {
	const e = "findOldestFile "
	dir_files, er := os.ReadDir(dir)
	if er != nil {
		err = e + er.Error()
		return
	}

	total_files = len(dir_files)

	oldestTime := time.Now()
	for _, dir_file := range dir_files {

		file, er := dir_file.Info()
		if er != nil {
			err = e + er.Error()
			return
		}

		if file.Mode().IsRegular() && file.ModTime().Before(oldestTime) {
			oldestFile = file
			oldestTime = file.ModTime()
		}
	}

	if oldestFile == nil {
		err = e + os.ErrNotExist.Error()
	}
	return
}
