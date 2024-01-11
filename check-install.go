package postgre

import (
	"fmt"
	"log"

	"github.com/cdvelop/objectdb"
)

// CHECK verifica el estado de la base motor PG POSTGRES
func (p *PG) CHECK(db *objectdb.Connection) bool {

	fmt.Printf(">>> *** Inicio Verificación de Base de datos [%v] ..... *** <<<\n\n", p.DBName)
	if ready := p.ExistDataBase(p.DBName, db); ready { //existe PG app?
		return true
	} else {
		//obtengo la conexión actual
		bkpDns := p.SwapConnect()
		// cierro la conexión
		defer db.Close()

		if ready = p.ExistDataBase("postgres", db); ready { //existe postgres?

			if ready = ExistDataBaseROL(bkpDns.UserDB, db); !ready { //verificar rol usuario app
				//crear rol app
				if err := db.QueryWithoutANSWER(`CREATE USER ` + bkpDns.UserDB + ` PASSWORD '` + p.passwordDB + `';`); err != "" {
					showErrorAndExit(fmt.Sprintf("!!! error al crear rol: %v", bkpDns.UserDB))
				}
				fmt.Printf(">>> rol: %v creado\n", bkpDns.UserDB)
			}

			err := db.QueryWithoutANSWER(`CREATE DATABASE ` + bkpDns.DBName + ` OWNER=` + bkpDns.UserDB + `;`)
			if err != "" {
				log.Fatal("error al crear base de datos:"+bkpDns.DBName, err)
			}

			fmt.Printf(">>> base de datos: %v creada\n", bkpDns.DBName)

			// cambiar clave administrador postgres****
			err = db.QueryWithoutANSWER(`ALTER USER postgres WITH PASSWORD '` + p.passwordDB + `';`)
			if err != "" {
				log.Fatal("error al cambiar password postgres", err)
			}
			fmt.Print(">>> password postgres actualizada\n")

			err = db.QueryWithoutANSWER(`GRANT ALL PRIVILEGES ON DATABASE ` + bkpDns.DBName + ` TO ` + bkpDns.UserDB + `;`)
			if err != "" {
				log.Fatalln("error al otorgar privilegios ", err)
			}

			fmt.Print(">>> privilegios rol " + bkpDns.UserDB + " otorgados a PG\n")

			//cambiar a nuevo dns inicial
			*p = bkpDns
			//seteo el pool de conexiones
			db.Set(p)
			fmt.Printf(">>> *** Verificación de Base de datos [%v] Finalizada *** <<<\n\n", p.DBName)
			return true

		} else {
			showErrorAndExit(fmt.Sprintf("!!! error motor de base de datos %v no existe chequear instalación", p.DataBasEngine()))
		}
	}
	return false
}
