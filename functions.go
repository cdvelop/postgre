package postgre

import (
	"fmt"
)

// SQLParameters string  //ej: postgres: "$" sqlite "?",
func (d *PG) PlaceHolders(index ...uint8) string {

	if len(index) == 1 {

		return fmt.Sprintf(`$%v`, index[0])
	} else {

		return `$`
	}

}

// SetListSyntax "%v=$%v"
func (d *PG) SetListSyntax(key string, i byte, set *[]string) {
	*set = append(*set, fmt.Sprintf("%v=$%v", key, i))
}

func (d *PG) TotalValuesSyntax(fields map[string]string) string {
	return fmt.Sprintf("$%v", len(fields))
}

func (d *PG) MakeSqInsertSyntax(i *byte, setValue *[]string) {
	*setValue = append(*setValue, fmt.Sprintf("$%v", *i+1))
	*i++
}

func (PG) DropTable() string {
	return "DROP TABLE IF EXISTS %v CASCADE;"
}

func (PG) SQLTableExist() string {
	return `SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = $1
	)`
}
