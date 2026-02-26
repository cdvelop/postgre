package postgre

import (
	"log"
	"os"

	"github.com/cdvelop/objectdb"
	_ "github.com/lib/pq"
)

// env_password_name ej: DB_PASSWORD_POSTGRE
func NewConnection(dba *PG, tables_in ...*struct {
	Name   string     //tabla ej: users,products,staff
	Legend string     // como se ve para el usuario ej: user por Usuarios
	Fields []struct { // campos
		Name   string // ej: id_user,name_user,phone
		Legend string // ej: "Nombre"
		Type   string // default TEXT
		Unique bool   //campo único e inalterable en db
	}
}) *objectdb.Connection {

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

	db := objectdb.Get(dba)

	dba.unixID = db

	// chequear tablas base de datos
	db.AddTablesToDB(tables_in, func(err string) {
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

			err = checkDir(dba.BackupDirectory)
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
