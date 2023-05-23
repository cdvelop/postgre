package postgre

import (
	"fmt"
	"sync"
	"time"
)

var (
	lastUnixIDatabase string
	controlProcess    sync.Mutex
)

// GetNewID retorna un id único para el ingreso a la base de datos tipo unix time
func (PG) GetNewID() string {
	controlProcess.Lock()
	idunix := fmt.Sprint(time.Now().UnixNano())
	for {
		// si no esta el id lo agregamos
		if lastUnixIDatabase != idunix { //last id unix time
			break
		} else {
			//obtenemos nueva marca
			idunix = fmt.Sprint(time.Now().UnixNano())
			// log.Printf(">>>new time id slow %v ", idunix)
		}
	}
	// log.Printf("unix time maps %v", setting.lastUnixIDatabase)
	lastUnixIDatabase = idunix //actualizo id
	controlProcess.Unlock()
	return idunix
}
