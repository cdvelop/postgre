# Implementation Plan: Postgres Adapter — ORM v0.0.13 DDL Support

## Development Rules
- **Single Responsibility Principle (SRP):** Every file must have a single, well-defined purpose.
- **Mandatory Dependency Injection (DI):** No global state. Interfaces for external dependencies.
- **Testing Runner (`gotest`):** Install first: `go install github.com/tinywasm/devflow/cmd/gotest@latest`
- **Standard Library Only in Tests:** NEVER use external assertion libraries.
- **Explicit Execution:** Do not modify source code unless all steps below are followed in order.

## Goal
Update `github.com/cdvelop/postgre` to support `tinywasm/orm@v0.0.13`. Breaking changes:
- `Model.Columns() []string` is replaced by `Model.Schema() []orm.Field`
- New DDL actions: `ActionCreateTable`, `ActionDropTable`, `ActionCreateDatabase` must be handled in `translate.go`
- All test models must implement the new `Schema()` interface
- FK constraints must generate standard PostgreSQL `CONSTRAINT fk_... FOREIGN KEY ... REFERENCES ...` clauses
- Integer PKs map to `SERIAL` or `BIGSERIAL`. Numeric autoincrement uses `SERIAL`

## References
- ORM Skill: `github.com/tinywasm/orm/docs/SKILL.md`
- ORM PLAN: `github.com/tinywasm/orm/docs/PLAN.md`

## Execution Steps

### Step 1 — Update dependency
```bash
go get github.com/tinywasm/orm@v0.1.0
go mod tidy
```

### Step 2 — Update `translate.go`
- **No change needed** for DML (still uses `q.Columns` populated by `db.go`).
- **Add DDL cases** in the main `translate` switch:

```go
case orm.ActionCreateTable:
    sb.Write("CREATE TABLE IF NOT EXISTS ")
    sb.Write(q.Table)
    sb.Write(" (")
    fields := m.Schema()
    for i, f := range fields {
        if i > 0 { sb.Write(", ") }
        sb.Write(f.Name)
        sb.Write(" ")
        isPK := f.Constraints&orm.ConstraintPK != 0
        isAuto := f.Constraints&orm.ConstraintAutoIncrement != 0
        if isPK && isAuto {
            // SERIAL/BIGSERIAL handles both PK and autoincrement in Postgres
            if f.Type == orm.TypeInt64 { sb.Write("BIGSERIAL") } else { sb.Write("SERIAL") }
        } else {
            sb.Write(postgresType(f.Type))
        }
        if isPK { sb.Write(" PRIMARY KEY") }
        if f.Constraints&orm.ConstraintNotNull != 0 { sb.Write(" NOT NULL") }
        if f.Constraints&orm.ConstraintUnique != 0 { sb.Write(" UNIQUE") }
    }
    // FK constraints as separate CONSTRAINT clauses (standard SQL)
    for _, f := range fields {
        if f.Ref != "" {
            refCol := f.RefColumn
            if refCol == "" { refCol = "id" }
            sb.Write(Sprintf(", CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(%s)",
                q.Table, f.Name, f.Name, f.Ref, refCol))
        }
    }
    sb.Write(")")

case orm.ActionDropTable:
    sb.Write("DROP TABLE IF EXISTS ")
    sb.Write(q.Table)

case orm.ActionCreateDatabase:
    sb.Write("CREATE DATABASE ")
    sb.Write(q.Database)
```

Add helper function:
```go
func postgresType(t orm.FieldType) string {
    switch t {
    case orm.TypeInt64:   return "BIGINT"
    case orm.TypeFloat64: return "DOUBLE PRECISION"
    case orm.TypeBool:    return "BOOLEAN"
    case orm.TypeBlob:    return "BYTEA"
    default:              return "TEXT"
    }
}
```

### Step 3 — Update `tests/adapter_test.go`
- Replace `Columns() []string` with `Schema() []orm.Field` on all test models (`User`).
- Add test for `CreateTable`:
```go
func TestCreateTable(t *testing.T) {
    // ... setup db
    err := dbORM.CreateTable(&User{})
    if err != nil { t.Errorf("CreateTable failed: %v", err) }
}
```
- Add test for `DropTable`.
- Add test for FK constraint generation via `Compile()` directly using `adapterClosed` or a mock.
- Add test for `CreateDatabase` action routing.
- Verify coverage includes all new DDL translate paths.

### Step 4 — Verify
```bash
gotest
```
Coverage must be ≥ 90%.

### Step 5 — Publish
```bash
gopush 'feat: support orm v0.0.13 DDL Schema API'
```
