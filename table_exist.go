package postgre

import (
	"log"

	"github.com/cdvelop/objectdb"
)

func (PG) TableExist(table_name string, db *objectdb.Connection) bool {

	filas, err := db.Query(`SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = $1
	)`, table_name)

	if err != nil {
		log.Println(err)
		return false
	}

	defer filas.Close()

	return filas.Next()
}
