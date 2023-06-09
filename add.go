package postgre

import (
	"fmt"
	"os"

	"github.com/cdvelop/dbtools"
	"github.com/cdvelop/model"
	"github.com/cdvelop/objectdb"
	_ "github.com/lib/pq"
)

// env_password_name ej: DB_PASSWORD_POSTGRE
func NewConnection(userDatabase, env_password_name, iPLocalServer, dataBasePORT, dataBaseName, directory_backup string, tables ...model.Object) *objectdb.Connection {

	password, existe := os.LookupEnv(env_password_name)
	if !existe {
		showErrorAndExit("valor password db variable de entorno: " + env_password_name + " no encontrado")
	}

	dba := PG{
		dataBaseName:     dataBaseName,
		ipLocalServer:    iPLocalServer,
		dataBasePORT:     dataBasePORT,
		userDatabase:     userDatabase,
		passwordDatabase: password,
		backup_directory: directory_backup,
		UnixID:           dbtools.NewUnixIdHandler(),
	}

	db := objectdb.Get(&dba)

	// chequear tablas base de datos
	for _, t := range tables {
		if !dba.TableExist(t.Name, db) {
			if !dbtools.CreateOneTABLE(db, t) {
				showErrorAndExit(fmt.Sprintf("no se logro crear tabla: %v", t.Name))
			}
		}
	}

	return db
}
