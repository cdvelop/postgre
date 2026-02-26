# PostgreSQL Adapter for tinywasm/orm

This repository implements the `orm.Adapter` interface for PostgreSQL, allowing it to be used with the `github.com/tinywasm/orm` library.

## Usage

```go
package main

import (
	"log"

	"github.com/cdvelop/postgre"
	"github.com/tinywasm/orm"
)

func main() {
	dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	adapter, err := postgre.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := orm.New(adapter)
	// Use db...
}
```

## Features

- Full `orm.Adapter` implementation.
- Transaction support via `BeginTx`.
- Secure SQL generation with parameterized queries using `$1`, `$2`, etc.
- Support for `Create`, `ReadOne`, `ReadAll`, `Update`, `Delete`.
- Efficient row scanning.

## Testing

Run tests using `go test ./tests/...`. Note that integration tests require a running PostgreSQL instance and `POSTGRES_DSN` environment variable set.
