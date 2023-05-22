package go_db_postgre_test

import (
	"sync"
	"testing"

	"github.com/cdvelop/postgre"
)

func Test_GetNewID(t *testing.T) {
	idRequired := 30
	wg := sync.WaitGroup{}
	wg.Add(idRequired)

	// PG := sqlite.NewConnection("test.PG", false)

	PG := postgre.NewConnection("test", "1", "127.0.0.1", "5432", "test", "./backup_test/")

	idObtained := make(map[string]int)
	var esperar sync.Mutex

	for i := 0; i < idRequired; i++ {
		go func() {
			defer wg.Done()
			id := PG.GetNewID()
			esperar.Lock()
			if cantId, exist := idObtained[id]; exist {
				idObtained[id] = cantId + 1
			} else {
				idObtained[id] = 1
			}
			esperar.Unlock()

		}()
	}
	wg.Wait()

	// fmt.Printf("total id requeridos: %v ob: %v\n", idRequired, len(idObtained))
	// fmt.Printf("%v", idObtained)
	if idRequired != len(idObtained) {
		t.Fail()
	}
}
