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

// New creates a new PostgresAdapter wrapped in an orm.DB.
// The dataSourceName parameter expects a valid PostgreSQL connection string.
// Example: "postgres://user:password@host:port/dbname?sslmode=disable"
func New(dataSourceName string) (*orm.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err // This branch might be hard to hit because sql.Open mostly validates the driver name.
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return orm.New(&PostgresAdapter{db: db}), nil
}

// AdapterForTest returns the raw PostgresAdapter for testing purposes.
func AdapterForTest(db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

// Ensure PostgresAdapter satisfies orm.Adapter.
var _ orm.Adapter = (*PostgresAdapter)(nil)
