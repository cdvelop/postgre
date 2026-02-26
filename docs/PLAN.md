# PostgreSQL Adapter - Phase 2 (Refinements)

This master prompt continues from the previous integration round, refining the API according to the latest domain requirements.

## Development Rules
- Constraints remain active: `gotest`, SRP, standard DI, pure stdlib testing, 500 lines limit, and flat hierarchy.
- **WASM / TinyGo Compatibility (`gotest` requirement):** You MUST replace all usages of standard `fmt` and `strings` with `github.com/tinywasm/fmt`. This is strictly required for the `gotest` command to successfully pass its validations.

## Execution Steps

### 1. Module Renaming
- Update `go.mod` to establish the definitive new ecosystem package path: `module github.com/tinywasm/postgres` (replacing the old `github.com/cdvelop/postgre`).
- Ensure any internal package references or internal files are updated natively to reflect `postgres` instead of `postgre` where applicable.

### 2. Direct ORM Injection (`adapter.go`)
- The user expressed that having to manually wrap `postgres.New(dsn)` with `orm.New()` is tedious.
- Refactor the constructor `New(dataSourceName string)` to **directly return an `*orm.DB`**.
- Example target signature: `func New(dataSourceName string) (*orm.DB, error)`.
- The internal logic of `New` will naturally create the internal `*sql.DB` connection, instantiate the `PostgresAdapter`, wrap it by calling `orm.New(adapter)`, and return the ready-to-use ORM `*DB` wrapper.
- The `PostgresAdapter` struct must remain internally mapped to safely satisfy `orm.Adapter` and `orm.TxAdapter`.

### 3. Tests & Verification (`tests/`)
- Ensure all tests in `tests/adapter_test.go` reflect the new constructor signature (they will receive `*orm.DB` directly and therefore have immediate access to `db.Create()`, `db.Tx()`, etc.).
- Validate implementation fully with `gotest`.
