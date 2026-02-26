package postgre

import (
	"context"
	"database/sql"

	"github.com/tinywasm/orm"
)

// PostgresTx wraps a SQL transaction and implements orm.TxBound.
type PostgresTx struct {
	tx *sql.Tx
}

// Ensure PostgresTx satisfies orm.TxBound.
var _ orm.TxBound = (*PostgresTx)(nil)

// Execute executes a query within the transaction.
func (p *PostgresTx) Execute(q orm.Query, m orm.Model, factory func() orm.Model, each func(orm.Model)) error {
	// Need to call `executeInternal` which takes an Executor.
	// sql.Tx implements Executor.
	return executeInternal(p.tx, q, m, factory, each)
}

// Commit commits the transaction.
func (p *PostgresTx) Commit() error {
	return p.tx.Commit()
}

// Rollback aborts the transaction.
func (p *PostgresTx) Rollback() error {
	return p.tx.Rollback()
}

// BeginTx starts a new transaction on the adapter.
// This overrides the placeholder in adapter.go (if we had kept it separate, but since `BeginTx` is part of `PostgresAdapter` methods).
// We'll define `BeginTx` here for `PostgresAdapter`.

// BeginTx starts a transaction.
func (p *PostgresAdapter) BeginTx() (orm.TxBound, error) {
	tx, err := p.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &PostgresTx{tx: tx}, nil
}
