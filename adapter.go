package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tinywasm/orm"
)

// PostgresAdapter is a PostgreSQL adapter for the tinywasm/orm library.
type PostgresAdapter struct {
	db *sql.DB
}

// New creates a new PostgresAdapter wrapped in an ORM DB.
func New(dataSourceName string) (*orm.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	adapter := &PostgresAdapter{db: db}
	return orm.New(adapter), nil
}

// Ensure PostgresAdapter satisfies orm.Adapter.
var _ orm.Adapter = (*PostgresAdapter)(nil)
