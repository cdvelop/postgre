package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/orm"
)

// Executor interface abstracts sql.DB and sql.Tx
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// executeInternal handles the actual execution logic for both DB and Tx
func executeInternal(exec Executor, q orm.Query, m orm.Model, factory func() orm.Model, each func(orm.Model)) error {
	queryStr, args, err := translate(q)
	if err != nil {
		return err
	}

	ctx := context.Background()

	switch q.Action {
	case orm.ActionCreate:
		// Attempt to use QueryRow to fetch ID if RETURNING id is supported/expected.
		// Since we don't know the PK column name generically, let's check if the query string
		// has `RETURNING` appended.
		// If not, just execute.

		// If `q` implies returning something (e.g., ORM sets flag), `translate` would have added it.
		// Currently `translate` does NOT add `RETURNING`.

		// Let's assume standard Exec for now.
		res, err := exec.ExecContext(ctx, queryStr, args...)
		if err != nil {
			return err
		}
		_ = res // We can't get LastInsertId from Postgres easily here without `RETURNING`.
		return nil

	case orm.ActionReadOne:
		row := exec.QueryRowContext(ctx, queryStr, args...)

		// m is the target model.
		dest := m.Pointers()

		// Scan into model
		if err := row.Scan(dest...); err != nil {
			if err == sql.ErrNoRows {
				return orm.ErrNotFound
			}
			return err
		}

		// The `each` callback is for iteration, but for ReadOne,
		// some adapters might call it once with the result.
		// The interface doc doesn't strictly say `each` is only for ReadAll.
		// However, `readOne` typically fills `m` directly.
		// If `each` is provided, we should call it?
		// Usually `each` is for `ReadAll`. `ReadOne` fills `m`.

		return nil

	case orm.ActionReadAll:
		rows, err := exec.QueryContext(ctx, queryStr, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			// Create new instance for each row
			newModel := factory()
			dest := newModel.Pointers()

			if err := rows.Scan(dest...); err != nil {
				return err
			}
			each(newModel)
		}
		return rows.Err()

	case orm.ActionUpdate, orm.ActionDelete:
		_, err := exec.ExecContext(ctx, queryStr, args...)
		return err

	default:
		return errors.New(fmt.Sprintf("unsupported action: %d", q.Action))
	}
}

// Execute implements orm.Adapter
func (p *PostgresAdapter) Execute(q orm.Query, m orm.Model, factory func() orm.Model, each func(orm.Model)) error {
	// We need to pass p.db (which is *sql.DB) but we need to satisfy Executor interface.
	// *sql.DB satisfies Executor.
	// But `Execute` signature on Adapter is `Execute(q, m, factory, each)`.
	// Our `executeInternal` takes `Executor` first.

	// Create context if needed, or use Background.
	return executeInternal(p.db, q, m, factory, each)
}
