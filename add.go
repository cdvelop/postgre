package postgre

import (
	_ "github.com/lib/pq"
)

func NewConnection(userDatabase, passwordDatabase, iPLocalServer, dataBasePORT, dataBaseName, directory_backup string) *PG {
	d := PG{
		dataBaseName:     dataBaseName,
		ipLocalServer:    iPLocalServer,
		dataBasePORT:     dataBasePORT,
		userDatabase:     userDatabase,
		passwordDatabase: passwordDatabase,
		backup_directory: directory_backup,
	}

	return &d
}
