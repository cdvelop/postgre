package postgre

import (
	"fmt"
	"log"

	"github.com/cdvelop/objectdb"
)

// CHECK verifica el estado de la base motor PG POSTGRES
func (p *PG) CHECK(db *objectdb.Connection) bool {

	fmt.Printf(">>> *** Inicio Verificación de Base de datos [%v] ..... *** <<<\n\n", p.dataBaseName)
	if ready := p.ExistDataBase(p.dataBaseName, db); ready { //existe PG app?
		return true
	} else {
		//obtengo la conexión actual
		bkpDns := p.SwapConnect()
		// cierro la conexión
		defer db.Close()

		if ready = p.ExistDataBase("postgres", db); ready { //existe postgres?

			if ready = p.ExistDataBaseROL(bkpDns.userDatabase, db); !ready { //verificar rol usuario app
				//crear rol app
				if ready = db.QueryWithoutANSWER(`CREATE USER `+bkpDns.userDatabase+` PASSWORD '`+p.passwordDatabase+`';`, ">>> creando rol PG"); !ready {
					log.Fatalf("!!! error al crear rol: %v", bkpDns.userDatabase)
				}
				fmt.Printf(">>> rol: %v creado\n", bkpDns.userDatabase)
			}

			if db.QueryWithoutANSWER(`CREATE DATABASE `+bkpDns.dataBaseName+` OWNER=`+bkpDns.userDatabase+`;`, ">>> creando PG") {
				fmt.Printf(">>> base de datos: %v creada\n", bkpDns.dataBaseName)
				// cambiar clave administrador postgres****
				if db.QueryWithoutANSWER(`ALTER USER postgres WITH PASSWORD '`+p.passwordDatabase+`';`, ">>> cambiando password PG") {
					fmt.Print(">>> password postgres actualizada\n")
				}
				if db.QueryWithoutANSWER(`GRANT ALL PRIVILEGES ON DATABASE `+bkpDns.dataBaseName+` TO `+bkpDns.userDatabase+`;`, ">>> otorgando privilegios a PG") {
					fmt.Print(">>> privilegios rol " + bkpDns.userDatabase + " otorgados a PG\n")

					//cambiar a nuevo dns inicial
					*p = bkpDns
					//seteo el pool de conexiones
					db.Set(p)
					fmt.Printf(">>> *** Verificación de Base de datos [%v] Finalizada *** <<<\n\n", p.dataBaseName)
					return true
				}
			}

		} else {
			log.Fatalf("!!! error motor de base de datos %v no existe chequear instalación", p.DataBasEngine())
		}
	}
	return false
}
