# Development

## Prerequisites

- Go 1.25 or later

## Commands

| Command | Description |
|---------|-------------|
| `go build ./...` | Build all packages |
| `go test ./...` | Run tests |
| `go test -v ./...` | Run tests with verbose output |
| `go test -v -run TestName ./path/to/package` | Run a single test |
| `go vet ./...` | Lint |
| `gofmt -w .` | Format code |
| `go mod tidy` | Tidy dependencies |

## Adding a Feature

New features follow the DDD layer structure:

1. **Domain** — Define the entity in `internal/domain/entities/`, add a repository interface in `internal/domain/repositories/`
2. **Application** — Create command structs in `internal/application/command/` and a service in `internal/application/services/`
3. **Infrastructure** — Implement the repository in `internal/infrastructure/db/sqlite/` and add a migration in `migrations/`
4. **Interface** — Add HTTP handlers in `internal/interface/web/handlers/` and register routes in `server.go`

## Code Conventions

- Follow standard Go conventions ([Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments))
- Use `gofmt` formatting
- Return errors, don't panic — wrap with `fmt.Errorf("context: %w", err)`
- Naming: `camelCase` for locals, `PascalCase` for exports, short receiver names
- Write table-driven tests
- Use `context.Context` as the first parameter where appropriate

## Database Migrations

Migrations live in `migrations/` and are embedded into the binary at build time. They run automatically on server startup.

To add a new migration, create up and down SQL files:

```
migrations/000002_add_feature.up.sql
migrations/000002_add_feature.down.sql
```

The migration runner uses [golang-migrate](https://github.com/golang-migrate/migrate).
