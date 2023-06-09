package postgre

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cdvelop/objectdb"
)

func (p *PG) DeleteAllTables(db *objectdb.Connection) {
	sql := `DROP SCHEMA IF EXISTS public CASCADE;
	        CREATE SCHEMA public;`
	if !db.QueryWithoutANSWER(sql, "") {
		log.Fatalln("ERROR EN DeleteAllTables")
	}
}

func (p *PG) DeleteDataBase() {

	PG, err := sql.Open("postgres", ConnectionString(p.passwordDatabase))
	if err != nil {
		log.Fatalf("¡Error de conexión %v!", err)
	}
	defer PG.Close()

	_, err = PG.Exec("DROP DATABASE IF EXISTS " + p.dataBaseName + ";")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("¡¡¡ BASE DE DATOS %v HA SIDO ELIMINADA!!!", p.dataBaseName)
}
