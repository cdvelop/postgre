package postgre_test

import (
	"log"
	"os"
	"testing"

	"github.com/cdvelop/postgre"
)

func Test_Postgres(t *testing.T) {

	const env_password_name = "DB_PASSWORD_POSTGRE_TEST"

	err := os.Setenv(env_password_name, "1")
	if err != nil {
		log.Fatal("No se logro setear variable de entorno ", env_password_name, err)
	}

	db := postgre.NewConnection(&postgre.PG{
		UserDB:          "test",
		EnvPasswordName: env_password_name,
		IPLocalServer:   "127.0.0.1",
		DataBasePORT:    "5432",
		DBName:          "test_db",
		BackupDirectory: "./backup_test/",
	})

	db.TestCrudStart(t)

	// Eliminar la variable de entorno
	err = os.Unsetenv(env_password_name)
	if err != nil {
		log.Fatalf("Error al remover la variable de entorno: %v %s\n", env_password_name, err)

	}
}
