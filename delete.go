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

// DeleteTABLE elimina tabla de una base de datos
func (p *PG) DeleteTABLE(table_name string, db *objectdb.Connection) {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %v CASCADE;", table_name)
	if _, err := db.Exec(sql); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">>> Tabla %v eliminada\n", table_name)
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
