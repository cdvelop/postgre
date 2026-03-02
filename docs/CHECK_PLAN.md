# Plan: Postgres Adapter Enhancement

## Development Rules
- **Prerequisites:** External agents must install `gotest` first:
  ```bash
  go install github.com/tinywasm/devflow/cmd/gotest@latest
  ```
- **Standard Library Only:** NEVER use external assertion libraries (e.g., testify, gomega). Use only the standard testing, net/http/httptest, and reflect APIs.
- **Testing Runner (gotest):** For Go tests, ALWAYS use the globally installed `gotest` CLI command. DO NOT use `go test` directly. Simply type `gotest` (no arguments) for the full suite.
- **WASM Compatibility:** Use `tinywasm/fmt` instead of `fmt`/`strings`/`strconv`/`errors`.
- **Single Responsibility Principle:** Every file must have a single purpose.

## Goal
Implement and verify `ActionCreateTable` and `ActionDropTable` in the Postgres adapter to support the `tinywasm/orm` DDL API, ensuring parity with the SQLite adapter and resolving potential integration gaps.

## Proposed Changes

### [Component] Tests
#### [NEW] [ddl_test.go](tests/ddl_test.go)
- Use minimal models for DDL verification:
  ```go
  type DDLModel struct {
      ID   int    `db:"pk,autoincrement"`
      Name string `db:"unique,not_null"`
  }
  ```
- Test `CreateTable` with various constraints: `ConstraintPK`, `ConstraintUnique`, `ConstraintNotNull`.
- Test auto-incrementing fields (`SERIAL`, `BIGSERIAL`).
- Test Foreign Key generation using a related minimal model.
- Verify `IF NOT EXISTS` behavior for `CreateTable` and `DropTable`.

### [Component] SQL Generation
#### [MODIFY] [translate.go](translate.go)
- Ensure `translate` correctly handles `orm.ActionCreateTable` and `orm.ActionDropTable`.
- Verify mapping of `orm.TypeInt64` + `PK` + `Auto` to `BIGSERIAL`.

## Verification Plan

### Automated Tests
- Run the full suite using `gotest` (requires `POSTGRES_DSN`):
```bash
gotest
```
