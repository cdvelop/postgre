package postgre

import (
	"os"
	"sync"

	"github.com/cdvelop/model"
	"github.com/cdvelop/objectdb"
	"github.com/cdvelop/timeserver"
	"github.com/cdvelop/unixid"
	_ "github.com/lib/pq"
)

// env_password_name ej: DB_PASSWORD_POSTGRE
func NewConnection(userDatabase, env_password_name, iPLocalServer, dataBasePORT, dataBaseName, directory_backup string, tables ...*model.Object) *objectdb.Connection {

	password, existe := os.LookupEnv(env_password_name)
	if !existe {
		showErrorAndExit("valor password db variable de entorno: " + env_password_name + " no encontrado")
	}

	uid, err := unixid.NewHandler(timeserver.Add(), &sync.Mutex{}, nil)
	if err != "" {
		showErrorAndExit(err)
	}

	dba := PG{
		dataBaseName:     dataBaseName,
		ipLocalServer:    iPLocalServer,
		dataBasePORT:     dataBasePORT,
		userDatabase:     userDatabase,
		passwordDatabase: password,
		backup_directory: directory_backup,
		UnixID:           uid,
	}

	db := objectdb.Get(&dba)

	// chequear tablas base de datos
	db.CreateTablesInDB(tables, func(err string) {
		if err != "" {
			showErrorAndExit(err)
		}
	})

	return db
}
