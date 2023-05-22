package postgre

// donde esta instalado postgres
var installation_directory string

func init() {
	// log.Println(">>> CHEQUEANDO DIRECTORIO INSTALACIÓN POSTGRES ...")
	//1 obtener ruta instalación postgres leer directorio c: por defecto
	pg_directory, _ := postgresWindowsDirectory()
	// log.Printf(">>> Postgres Version %v Directorio Instalación: \n%v\n", ver, pg_directory)
	installation_directory = pg_directory
}
