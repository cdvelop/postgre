package postgre

import (
	"fmt"
	"os"
)

func postgresWindowsDirectory() (out_path, version string) {
	const default_installation_path = "%v:/Program Files/PostgreSQL"
	possible_letters := []string{"C", "D", "F"}

	for _, letter := range possible_letters {
		temp_path := fmt.Sprintf(default_installation_path, letter)

		files, err := os.ReadDir(temp_path)
		if err == nil {
			for _, f := range files {
				if f.Name() != "" {
					version = f.Name()
					out_path = temp_path + "/" + f.Name()
					break
				}
			}
		}
	}

	if out_path == "" {
		showErrorAndExit("¡¡¡ POSTGRES SQL NO INSTALADO !!!")
	}

	return
}
