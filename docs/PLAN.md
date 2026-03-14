# Migrate PostgreSQL Adapter to fmt.Field API

## Context

The ORM now uses `fmt.Field` (from `tinywasm/fmt`) with individual bool constraint fields instead of `orm.Field` with bitmask constraints. The PostgreSQL adapter's `translate()` function reads field metadata to build SQL — it must be updated to read the new `fmt.Field` struct.

### Key API Changes

| Old (current) | New (target) |
|---|---|
| `orm.Field` with `Constraints` bitmask | `fmt.Field` with `PK`, `Unique`, `NotNull`, `AutoInc` bools |
| `orm.TypeText`, `orm.TypeInt64`, `orm.TypeFloat64`, `orm.TypeBool`, `orm.TypeBlob` | `fmt.FieldText`, `fmt.FieldInt`, `fmt.FieldFloat`, `fmt.FieldBool`, `fmt.FieldBlob` |
| `f.Constraints&orm.ConstraintPK != 0` | `f.PK` |
| `f.Constraints&orm.ConstraintUnique != 0` | `f.Unique` |
| `f.Constraints&orm.ConstraintNotNull != 0` | `f.NotNull` |
| `f.Constraints&orm.ConstraintAutoIncrement != 0` | `f.AutoInc` |
| `m.Values()` | `fmt.ReadValues(m.Schema(), m.Pointers())` |
| `orm.FieldExt{orm.Field, Ref, RefColumn}` | `orm.FieldExt{fmt.Field, Ref, RefColumn}` |

### Target fmt.Field Struct (`tinywasm/fmt`)

```go
type Field struct {
    Name    string
    Type    FieldType // FieldText, FieldInt, FieldFloat, FieldBool, FieldBlob, FieldStruct
    PK      bool
    Unique  bool
    NotNull bool
    AutoInc bool
    Input   string // UI hint for form layer
    JSON    string // JSON key ("email,omitempty"). Empty = use Field.Name
}
```

### FieldExt (FK metadata, used by adapters)

```go
type FieldExt struct {
    fmt.Field
    Ref       string // FK: target table name. Empty = no FK.
    RefColumn string // FK: target column. Empty = auto-detect PK.
}
```

### ORM Model Interface (new)

```go
type Model interface {
    fmt.Fielder           // Schema() []fmt.Field + Pointers() []any
    TableName() string
}
```

Values are obtained via `fmt.ReadValues(m.Schema(), m.Pointers())` — no more `m.Values()` method.

---

## Stage 1 — Update translate.go

**File**: `translate.go`

1. Update `postgresType()` function:
   - Replace `orm.TypeInt64` → `fmt.FieldInt`
   - Replace `orm.TypeFloat64` → `fmt.FieldFloat`
   - Replace `orm.TypeBool` → `fmt.FieldBool`
   - Replace `orm.TypeBlob` → `fmt.FieldBlob`
   - Default case (text) remains for `fmt.FieldText`

2. Update `ActionCreateTable` handling in `translate()`:
   - Replace all bitmask constraint checks with bool field access:
     - `f.Constraints&orm.ConstraintPK != 0` → `f.PK`
     - `f.Constraints&orm.ConstraintUnique != 0` → `f.Unique`
     - `f.Constraints&orm.ConstraintNotNull != 0` → `f.NotNull`
     - `f.Constraints&orm.ConstraintAutoIncrement != 0` → `f.AutoInc`
   - FK metadata: `orm.FieldExt` now embeds `fmt.Field` — access `ext.Ref` and `ext.RefColumn` unchanged

3. Update `ActionCreate` and `ActionUpdate` handling:
   - If code calls `m.Values()`, replace with `fmt.ReadValues(m.Schema(), m.Pointers())`

4. Add `"github.com/tinywasm/fmt"` import

---

## Stage 2 — Update adapter.go

**File**: `adapter.go`

1. If `Compile()` references `orm.Field` directly, update to `fmt.Field`
2. Verify `orm.Model` interface still works (it now embeds `fmt.Fielder`)

---

## Stage 3 — Update Tests

**Files**: `export_test.go`, `postgres_translate_test.go`

1. Update any test fixtures that construct `orm.Field` literals → `fmt.Field`
2. Update bitmask constraint assertions → bool field assertions
3. Update type constant references

---

## Stage 4 — Update go.mod

1. Run `go mod tidy`
2. Ensure `tinywasm/fmt` is at latest version

---

## Verification

```bash
gotest
```

## Linked Documents

- [POSTGRES_SETUP.md](POSTGRES_SETUP.md)
