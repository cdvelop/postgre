# PostgreSQL Adapter Implementation

This document is the **Master Prompt (PLAN.md)** for refactoring the `postgre` library to act as an Adapter for the `tinywasm/orm` ecosystem. Every execution agent must follow this plan sequentially.

---

## Development Rules

- **SRP (Single Responsibility Principle):** Every file must have a single, well-defined purpose.
- **Mandatory Dependency Injection (DI):** The library acts as a dependency to be injected. It must import `github.com/tinywasm/orm`.
- **Flat Hierarchy:** Keep files in the root.
- **Max 500 lines per file:** Subdivide by domain if needed.
- **Test Organization:** Move tests to `tests/` if >5 test files exist.
- **Testing Runner:** Always use `gotest`.
- **WASM Compatibility:** Backend-only. `//go:build !wasm` where applicable.
- **No external assertion libs:** Standard `testing` package only.
- **Documentation First:** Update `README.md` concurrently or before source code files.

---

## Architecture Overview

Replace current custom CRUD operations with a structure implementing `orm.Adapter`:

```go
type PostgresAdapter struct {
    pool *pgxpool.Pool // or *sql.DB depending on current usage
}

func (p *PostgresAdapter) Execute(q orm.Query, m orm.Model, factory func() orm.Model, each func(orm.Model)) error
```

Transactions must implement `orm.TxAdapter` and `orm.TxBound`.

---

## Execution Phases

### Phase 1: Struct definition and Connection (`adapter.go`)
1. Define `PostgresAdapter` encapsulating the connection pool.
2. Update `go.mod` to require `github.com/tinywasm/orm`.
3. Assure connection creation yields a `*PostgresAdapter`.

### Phase 2: Translation Engine (`translate.go`)
1. Create a translator to map `orm.Query` to Postgres SQL dialect.
2. **Crucial Difference:** Use `$1`, `$2`, `$3` positional numeric placeholders rather than `?`. Maintain an ongoing index count when building `where` clauses and updates.
3. If necessary under the `orm` specifications, for `ActionCreate` and `ActionUpdate`, append `RETURNING id` if required by downstream handlers. The ORM strictly dictates data handling, handle Postgres specifics accordingly.

### Phase 3: Adapter Implementation (`execute.go`)
1. Implement the generic `Execute` mapping on `PostgresAdapter`.
2. Dispatch action subsets (`Create`, `ReadOne`, `ReadAll`, `Update`, `Delete`) to the pool backend.
3. Map `orm.ActionReadAll` to `db.Query`, loop via `Rows.Scan` explicitly into `m.Pointers()`, and invoke the callback via `each(m)`. 

### Phase 4: Transaction Support (`tx.go`)
1. Satisfy `orm.TxAdapter` providing `BeginTx() (orm.TxBound, error)` yielding a `PostgresTxBound` that tracks an active transaction block.
2. Reroute structural `Execute` requests implicitly executing within the bounded transaction.

### Phase 5: Cleanup and Testing
1. Erase obsolete custom APIs logic files (`add.go`, `delete.go`, `file_old.go`, `functions.go`, etc.).
2. Re-anchor tests integrating `gotest` with a temporary Postgres server configuration schema.
3. Formally verify parameterized numeric placeholder integrity.
