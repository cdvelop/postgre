package postgre

import (
	"log"

	"github.com/cdvelop/objectdb"
)

func (d *PG) DataBaseMaintenance(db *objectdb.Connection) {
	err := db.QueryWithoutANSWER(`VACUUM FULL ANALYZE;`)
	if err != "" {
		log.Println("error mantenimiento db: " + err)

	} else {
		log.Println("mantenimiento db:" + d.DBName + " ok")
	}
}
