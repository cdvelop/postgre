package postgre

import (
	"log"

	"github.com/cdvelop/objectdb"
)

func (d *PG) DataBaseMaintenance(db *objectdb.Connection) {
	res := db.QueryWithoutANSWER(`VACUUM FULL ANALYZE;`, "Realizando Mantenimiento Db")
	log.Printf("Mantenimiento DB Finalizado Correctamente? [%v]\n", res)
}
