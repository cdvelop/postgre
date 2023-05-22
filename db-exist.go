package postgre

import (
	"fmt"
	"log"

	"github.com/cdvelop/objectdb"
)

// ExistDataBase verifica si existe una base de datos postgres determinada
func (d *PG) ExistDataBase(data_base_name string, db *objectdb.Connection) bool {
	sql := fmt.Sprintf(d.DBExists(), data_base_name)
	// log.Printf("sql %v", sql)
	return db.Exists("base de datos", data_base_name, sql)
}

// DeleleDataBASE ..
func (p *PG) DeleleDataBASE(db *objectdb.Connection) {
	dnsBkp := p.SwapConnect()
	//seteo el pool de conexiones
	db.Set(p)

	defer db.Close()
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %v;", dnsBkp.dataBaseName)
	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("error %v al eliminar base de datos sql: %v ", err, sql)
	}
	*p = dnsBkp
	fmt.Printf(">>> Base de datos [%v] Eliminada !!!\n", dnsBkp.dataBaseName)
}
