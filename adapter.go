package postgre

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tinywasm/orm"
)

// PostgresAdapter is a PostgreSQL adapter for the tinywasm/orm library.
type PostgresAdapter struct {
	db *sql.DB
}

// New creates a new PostgresAdapter.
func New(dataSourceName string) (*PostgresAdapter, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresAdapter{db: db}, nil
}

// Ensure PostgresAdapter satisfies orm.Adapter.
var _ orm.Adapter = (*PostgresAdapter)(nil)
