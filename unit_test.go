package postgre_test

import (
	"os"
	"testing"

	"github.com/cdvelop/objectdb"
	"github.com/cdvelop/postgre"
)

func Test_Postgres(t *testing.T) {
	// test......
	user_db := "test"
	db_name := "test_db"

	password := os.Getenv("DB_PASSWORD_POSTGRE")
	if password == "" {
		password = "1"
	}

	// 1- conexión por defecto postgresql
	pg := postgre.NewConnection("postgres", password, "127.0.0.1", "5432", "postgres", "./backup_test/")
	db := objectdb.Get(pg)

	// 2- crear usuario si no existe
	if !pg.CreateUserRolDB(user_db, password, db) {
		return
	}

	// 3- setear conexión con datos actuales
	*pg = *postgre.NewConnection(user_db, password, "127.0.0.1", "5432", db_name, "./backup_test/")
	db.Set(db)

	// 4- chequear base de datos
	if !pg.CHECK(db) {
		return
	}

	// 5- eliminar eliminar las 3 tablas del modelo usado en objectdb
	for _, table := range []string{"usuario", "especialidad", "credentials"} {
		pg.DeleteTABLE(table, db)
	}

	db.TestCrudStart(t, pg)

	//test......
}
