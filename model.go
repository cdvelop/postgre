package postgre

// PG formato cadena conexión
type PG struct {
	dataBaseName     string //nombre de la base de datos
	ipLocalServer    string //ip servidor donde estará la base de datos
	dataBasePORT     string //puerto
	userDatabase     string //usuario base de datos
	passwordDatabase string //contraseña
	backup_directory string //ej "D:\postgres_backup"

}
