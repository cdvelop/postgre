package postgre

import "fmt"

//DataBasEngine "postgres"
func (d PG) DataBasEngine() string {
	return "postgres"
}

//ej "mydb"
func (d PG) DataBaseName() string {
	return d.DBName
}

// ConnectionString formato cadena de conexión
// postgres:// user : passwordDB  @  127.0.0.1  :  5432 / nombrebasedatos   ?sslmode=disable"
func (d *PG) ConnectionString() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", d.UserDB, d.passwordDB, d.IPLocalServer, d.DataBasePORT, d.DBName)
}

// conexión string por defecto postgres
func ConnectionString(password string) string {
	return fmt.Sprintf("postgres://postgres:%v@127.0.0.1:5432/postgres?sslmode=disable", password)
}
