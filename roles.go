package postgre

import (
	"fmt"
	"log"

	"github.com/cdvelop/objectdb"
)

func (p *PG) DeleteRolDataBase(rol_name string, db *objectdb.Connection) {
	dnsBkp := p.SwapConnect()
	db.Set(p) //seteo el pool de conexiones
	newdb := db
	defer newdb.Close()
	sql := fmt.Sprintf("DROP ROLE IF EXISTS %v;", dnsBkp.userDatabase)
	if _, err := newdb.Exec(sql); err != nil {
		log.Fatalf("error %v al eliminar rol usuario base de datos sql: %v ", err, sql)
	}
	*p = dnsBkp
	fmt.Printf(">>> Usuario PG Rol [%v] Eliminado !!!\n", dnsBkp.dataBaseName)
}

// ExistDataBaseROL verifica  si rol usuario aplicación existe
func (p *PG) ExistDataBaseROL(rol string, db *objectdb.Connection) bool {
	sql := fmt.Sprintf(p.ROLExists(), rol)

	return db.Exists("rol", rol, sql)
}

func (p *PG) CreateUserRolDB(user_name, password string, db *objectdb.Connection) bool {
	db.Set(p)
	defer db.Close()

	if ready := p.ExistDataBaseROL(user_name, db); !ready { //verificar rol usuario app

		//crear rol app
		if ready := db.QueryWithoutANSWER(`CREATE USER `+user_name+` PASSWORD '`+password+`';`, ">>> creando rol PG"); !ready {
			log.Fatalf("!!! error al crear rol: %v", user_name)
			return false
		}

	}

	// sql := fmt.Sprintf(p.ROLExists(), user_name)

	return true
}
