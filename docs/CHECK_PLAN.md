# PLAN — tinywasm/postgres: Document Update/Delete Condition Contract

## Context

A critical data-safety bug was found: `db.Update(&model)` without conditions
generates a full-table UPDATE with no WHERE clause.

**The fix lives entirely in `tinywasm/orm`** (see `tinywasm/orm/docs/PLAN.md`).
After that fix, `Update(m, cond, rest...)` requires at least one `Condition`
at compile time — it is impossible to call with zero conditions.

**This library requires no code changes.** The `ActionUpdate` case in `translate()`
already handles conditions correctly. Since `orm` now guarantees
`q.Conditions` is never empty for an Update query, the bug cannot reach this driver.

This plan exists to **update the documentation** and add a white-box regression
test verifying the driver correctly handles explicit conditions in UPDATE statements.

---

## Development Rules

- Standard Library only — no external dependencies.
- Max 500 lines per file.
- Run tests with `gotest` (never `go test` directly).
- Publish with `gopush 'message'` after all tests pass.
- **Documentation must be updated before touching code.**
- Prerequisites in isolated environments:
  ```bash
  go install github.com/tinywasm/devflow/cmd/gotest@latest
  ```

---

## Step 1 — Update `README.md`

In any usage examples showing `db.Update(...)`, ensure at least one explicit
condition is always present. Update any snippet that showed zero-arg Update.

Example of what to add to the README:

```markdown
## Update

`db.Update` always requires at least one `Condition`. This is enforced at
compile time by `tinywasm/orm`. There is no "update by PK implicitly" magic.

```go
// ✅ Correct
if err := db.Update(&user, orm.Eq(User_.ID, user.ID)); err != nil { ... }

// ❌ Compile error (caught by tinywasm/orm — will not reach the PostgreSQL layer)
db.Update(&user)
```
```

---

## Step 2 — White-box regression test (translate layer)

Verify the postgres `translate()` function correctly generates a `WHERE` clause
when conditions are provided. This confirms the driver is not vulnerable even
if the ORM layer were bypassed.

Create or update `tests/postgres_translate_test.go`:

```go
//go:build !wasm

package postgre_test

import (
	"strings"
	"testing"

	"github.com/tinywasm/orm"
	postgre "github.com/tinywasm/postgres"
)

// testUserModel is a minimal test model with a TEXT primary key (string PK, like unixid).
type testUserModel struct {
	ID   string
	Name string
	Age  int
}

func (u *testUserModel) TableName() string { return "users" }
func (u *testUserModel) Schema() []orm.Field {
	return []orm.Field{
		{Name: "id", Type: orm.TypeText, Constraints: orm.ConstraintPK},
		{Name: "name", Type: orm.TypeText},
		{Name: "age", Type: orm.TypeInt64},
	}
}
func (u *testUserModel) Values() []any   { return []any{u.ID, u.Name, u.Age} }
func (u *testUserModel) Pointers() []any { return []any{&u.ID, &u.Name, &u.Age} }

// TestTranslate_Update_WithCondition verifies that translate() generates a
// valid UPDATE ... SET ... WHERE ... when at least one condition is present.
// This is the contract guaranteed by tinywasm/orm's mandatory first Condition.
func TestTranslate_Update_WithCondition(t *testing.T) {
	m := &testUserModel{ID: "abc123", Name: "Alice", Age: 30}
	q := orm.Query{
		Action:  orm.ActionUpdate,
		Table:   "users",
		Columns: []string{"id", "name", "age"},
		Values:  m.Values(),
		// At least one condition — as guaranteed by tinywasm/orm after the fix.
		Conditions: []orm.Condition{orm.Eq("id", "abc123")},
	}

	sql, args, err := postgre.ExportTranslate(q, m)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Must contain WHERE clause.
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE clause in UPDATE, got: %s", sql)
	}

	// Must use parameterized form ($N for postgres).
	if !strings.Contains(sql, "$") {
		t.Errorf("expected parameterized query, got: %s", sql)
	}

	// Condition value must appear in args.
	pkFound := false
	for _, a := range args {
		if a == "abc123" {
			pkFound = true
			break
		}
	}
	if !pkFound {
		t.Errorf("expected PK value 'abc123' in args, got: %v", args)
	}
}

// TestTranslate_Update_MultipleConditions verifies that AND conditions in an
// UPDATE query produce correct SQL.
func TestTranslate_Update_MultipleConditions(t *testing.T) {
	m := &testUserModel{ID: "abc123", Name: "Alice", Age: 30}
	q := orm.Query{
		Action:  orm.ActionUpdate,
		Table:   "users",
		Columns: []string{"id", "name", "age"},
		Values:  m.Values(),
		Conditions: []orm.Condition{
			orm.Eq("id", "abc123"),
			orm.Eq("name", "Alice"),
		},
	}

	sql, args, err := postgre.ExportTranslate(q, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("expected WHERE in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "AND") {
		t.Errorf("expected AND between conditions, got: %s", sql)
	}
	_ = args
}
```

> **`ExportTranslate` helper:** If not already present, create `export_test.go`
> in the package root:
>
> ```go
> package postgre
>
> import "github.com/tinywasm/orm"
>
> // ExportTranslate exposes translate for white-box testing.
> func ExportTranslate(q orm.Query, m orm.Model) (string, []any, error) {
>     return translate(q, m)
> }
> ```

---

## Acceptance Criteria

- [ ] `README.md` updated — no zero-condition Update example
- [ ] `ExportTranslate` helper present in `export_test.go` (create if missing)
- [ ] `TestTranslate_Update_WithCondition` passes
- [ ] `TestTranslate_Update_MultipleConditions` passes
- [ ] All existing postgres tests still pass (`gotest`)
- [ ] `gopush 'docs+test: document Update condition contract and add translate regression tests'`
