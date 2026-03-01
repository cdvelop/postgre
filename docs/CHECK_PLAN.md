# Implementation Plan: Eliminate Global State from Postgres Adapter

## Development Rules
- **WASM Environment (`tinywasm`):** Frontend Go Compatibility requires standard library replacements (`tinywasm/fmt`).
- **Single Responsibility Principle (SRP):** Every file must have a single, well-defined purpose.
- **Mandatory Dependency Injection (DI):** No global state. Interfaces for external dependencies.
- **Testing Runner (`gotest`):** ALWAYS use the globally installed `gotest` CLI command.
- **Documentation First:** Update docs before coding.

## Goal
The `postgres` adapter currently uses a global registry (`dbRegistry` and `dbMu`) to associate an `*orm.DB` instance with its underlying driver connection. This was necessary because `*orm.DB` previously encapsulated its `Executor` privately without exposing a `Close()` or `RawExecutor()` method. Now that `github.com/tinywasm/orm` (v0.0.10) natively exposes these methods, the global state can be completely removed to comply with the strictly required DI and SRP rules.

## Execution Steps

### 1. Update `go.mod`
- Update the `github.com/tinywasm/orm` dependency to the latest version.
- Run `go get github.com/tinywasm/orm@v0.0.10` to pull the `Close()` and `RawExecutor()` methods.

### 2. Remove Global State from `adapter.go`
- Delete `dbRegistry` and `dbMu` variables.
- Update `postgre.New` to no longer register the database connection to `dbRegistry` and `dbMu`.

### 3. Update `Close` and `ExecSQL` inside `adapter.go`
- Refactor the `Close` function to directly call `db.Close()`.
- Refactor the `ExecSQL` function to retrieve the raw executor via `db.RawExecutor()`, cast it to a type that has `.Exec(...)`, or use it properly if it covers your requirements.
- **Note**: Ensure `tinywasm/fmt` is used consistently for the error handling if required (`Err` or `Errf`).

### 4. Verify Tests
- Run `gotest` with the correct `POSTGRES_DSN` in `.env`.
- Ensure all tests pass.
- Ensure test coverage remains >90%.

### 5. Update Documentation
- Check `README.md` and any API documentation to reflect the cleaner implementation architecture.

## Verification Plan
### Automated Tests
- Run `gotest` in `tinywasm/postgres` to verify that `Close()` and `ExecSQL()` continue behaving correctly without the global state registry.
