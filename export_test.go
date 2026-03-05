package postgre

import "github.com/tinywasm/orm"

// ExportTranslate exposes translate for white-box testing.
func ExportTranslate(q orm.Query, m orm.Model) (string, []any, error) {
	return translate(q, m)
}
