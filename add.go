package postgre

import (
	"log"
	"os"
	"sync"

	"github.com/cdvelop/model"
	"github.com/cdvelop/objectdb"
	"github.com/cdvelop/timeserver"
	"github.com/cdvelop/unixid"
	_ "github.com/lib/pq"
)

// env_password_name ej: DB_PASSWORD_POSTGRE
func NewConnection(dba *PG, tables ...*model.Object) *objectdb.Connection {

	const e = "postgres new connection error: "

	password, existe := os.LookupEnv(dba.EnvPasswordName)
	if !existe {
		showErrorAndExit("valor password db variable de entorno: " + dba.EnvPasswordName + " no encontrado")
	}

	if dba.DataBasePORT == "" {
		dba.DataBasePORT = "5432"
	}

	if dba.IPLocalServer == "" {
		dba.IPLocalServer = "127.0.0.1"
	}

	if dba.BackupDirectory == "" {
		dba.BackupDirectory = "./backupdb"
	}

	if dba.TotalBackupsMaintain == 0 {
		dba.TotalBackupsMaintain = 1
	}

	dba.passwordDB = password

	uid, err := unixid.NewHandler(timeserver.Add(), &sync.Mutex{}, unixid.NoSessionNumber{})
	if err != "" {
		showErrorAndExit(e + err)
	}

	dba.idUnix = uid

	db := objectdb.Get(dba)

	// chequear tablas base de datos
	db.CreateTablesInDB(tables, func(err string) {
		if err != "" {
			showErrorAndExit(e + err)
		}
	})

	if dba.ScheduleAdapter != nil {

		if dba.BackupDirectory != "" && dba.ScheduleBackup != "" {

			err := dba.AddFuncToSchedule(dba.ScheduleBackup, "respaldo base de datos:"+dba.DBName, dba.BackupDataBase)
			if err != "" {
				log.Println(e + err)
			}
		}

		if dba.ScheduleMaintenance != "" {

			err := dba.AddFuncToSchedule(dba.ScheduleMaintenance, "mantenimiento base de datos:"+dba.DBName, dba.DataBaseMaintenance, db)
			if err != "" {
				log.Println(e + err)
			}

		}

	}

	return db
}
