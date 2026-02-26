package tests

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/cdvelop/postgre"
	"github.com/tinywasm/orm"
)

// Define a simple model for testing
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Columns() []string {
	return []string{"id", "name", "email", "created_at"}
}

func (u *User) Values() []any {
	return []any{u.ID, u.Name, u.Email, u.CreatedAt}
}

func (u *User) Pointers() []any {
	return []any{&u.ID, &u.Name, &u.Email, &u.CreatedAt}
}

// Factory function
func NewUser() orm.Model {
	return &User{}
}

func TestPostgresAdapter(t *testing.T) {
	// Setup connection
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	}

	// Try to connect to DB for setup
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping test: could not connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: database not reachable: %v", err)
	}

	// Setup Schema
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT,
			created_at TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	// Initialize Adapter
	adapter, err := postgre.New(dsn)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	// Test Create
	user := &User{Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()}
	qCreate := orm.Query{
		Action:  orm.ActionCreate,
		Table:   user.TableName(),
		Columns: []string{"name", "email", "created_at"},
		Values:  []any{user.Name, user.Email, user.CreatedAt},
	}
	if err := adapter.Execute(qCreate, user, nil, nil); err != nil {
		t.Errorf("Create failed: %v", err)
	}

	// Test ReadAll
	qReadAll := orm.Query{
		Action: orm.ActionReadAll,
		Table:  user.TableName(),
	}
	var users []orm.Model
	err = adapter.Execute(qReadAll, nil, NewUser, func(m orm.Model) {
		users = append(users, m)
	})
	if err != nil {
		t.Errorf("ReadAll failed: %v", err)
	}
	if len(users) == 0 {
		t.Errorf("Expected users, got 0")
	}

	// Test ReadOne (by Name)
	qReadOne := orm.Query{
		Action: orm.ActionReadOne,
		Table:  user.TableName(),
		Conditions: []orm.Condition{orm.Eq("name", "Alice")},
	}
	foundUser := &User{}
	if err := adapter.Execute(qReadOne, foundUser, nil, nil); err != nil {
		t.Errorf("ReadOne failed: %v", err)
	}
	if foundUser.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", foundUser.Name)
	}

	// Test Update
	qUpdate := orm.Query{
		Action:     orm.ActionUpdate,
		Table:      user.TableName(),
		Columns:    []string{"email"},
		Values:     []any{"alice_new@example.com"},
		Conditions: []orm.Condition{orm.Eq("name", "Alice")},
	}
	if err := adapter.Execute(qUpdate, nil, nil, nil); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Test Delete
	qDelete := orm.Query{
		Action:     orm.ActionDelete,
		Table:      user.TableName(),
		Conditions: []orm.Condition{orm.Eq("name", "Alice")},
	}
	if err := adapter.Execute(qDelete, nil, nil, nil); err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}
