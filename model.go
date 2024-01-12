package postgre

import "github.com/cdvelop/unixid"

type ScheduleAdapter interface {
	AddFuncToSchedule(schedule, description string, fun any, args ...any) (err string)
}

// PG formato cadena conexión
type PG struct {
	DBName          string //nombre de la base de datos
	UserDB          string //usuario base de datos
	EnvPasswordName string // ej: APP_PASSWORD

	IPLocalServer string //ip servidor donde estará la base de datos default: "127.0.0.1"
	DataBasePORT  string //puerto default: 5432

	BackupDirectory      string //ej "D:\postgres_backup" default: ./backupdb
	ScheduleBackup       string //ej: "0 14,19 * * 1-6" = días laborales a las 14:00 y 19:00 se realizara respaldo
	TotalBackupsMaintain int    // total de respaldos a mantener por defecto solo 1
	ScheduleMaintenance  string // ej: "0 0 1 * 0" = primer domingo de cada mes a las 00:00  se realizara mantenimiento
	ScheduleAdapter             // adaptador tareas programadas (cronograma) = AddFuncToSchedule(schedule, description string, fun any, args ...any) (err string)

	passwordDB string //contraseña

	idUnix *unixid.UnixID
}
