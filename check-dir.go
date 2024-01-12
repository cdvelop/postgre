package postgre

import "os"

func checkDir(dir string) (err string) {

	_, er := os.Stat(dir)

	if os.IsNotExist(er) {
		er := os.MkdirAll(dir, os.ModePerm)
		if er != nil {
			return "checkCertDir al crear el directorio: " + er.Error()
		}
	}

	return ""
}
