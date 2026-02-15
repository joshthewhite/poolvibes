# Development

## Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev) — `go install github.com/go-task/task/v3/cmd/task@latest`
- [templ](https://templ.guide) CLI — `go install github.com/a-h/templ/cmd/templ@latest` (only needed when editing `.templ` files)
- [air](https://github.com/air-verse/air) — `go install github.com/air-verse/air@latest` (only needed for `task dev`)

## Commands

| Command | Description |
|---------|-------------|
| `task --list` | Show all available tasks |
| `task build` | Build binary to `bin/poolvibes` |
| `task test` | Run tests |
| `task test:verbose` | Run tests with verbose output |
| `go test -v -run TestName ./path/to/package` | Run a single test |
| `task lint` | Lint |
| `task fmt` | Format code |
| `task templ` | Regenerate Go code from `.templ` files |
| `task dev` | Start dev server with live reload (templ + air) |
| `task run` | Build and run the server |
| `task clean` | Remove build artifacts |
| `task tidy` | Tidy dependencies |
| `task docker:build` | Build Docker image |
| `task docker:up` | Start Docker services |
| `task docker:down` | Stop Docker services |

## Adding a Feature

New features follow the DDD layer structure:

1. **Domain** — Define the entity in `internal/domain/entities/`, add a repository interface in `internal/domain/repositories/`
2. **Application** — Create command structs in `internal/application/command/` and a service in `internal/application/services/`
3. **Infrastructure** — Implement the repository in `internal/infrastructure/db/sqlite/` and add a migration in `migrations/`
4. **Interface** — Add templ components in `internal/interface/web/templates/`, HTTP handlers in `internal/interface/web/handlers/`, and register routes in `server.go`

## Code Conventions

- Follow standard Go conventions ([Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments))
- Use `gofmt` formatting
- Return errors, don't panic — wrap with `fmt.Errorf("context: %w", err)`
- Naming: `camelCase` for locals, `PascalCase` for exports, short receiver names
- Write table-driven tests
- Use `context.Context` as the first parameter where appropriate

## CI

GitHub Actions runs on every push to `main` and on pull requests. The workflow (`.github/workflows/ci.yml`) runs:

1. `go build ./...`
2. `go vet ./...`
3. `go test ./...`

Railway's **Wait for CI** integration ensures deployments only proceed after CI passes.

## Database Migrations

Migrations live in `migrations/` and are embedded into the binary at build time. They run automatically on server startup.

To add a new migration, create up and down SQL files:

```
migrations/000002_add_feature.up.sql
migrations/000002_add_feature.down.sql
```

The migration runner uses [golang-migrate](https://github.com/golang-migrate/migrate).
