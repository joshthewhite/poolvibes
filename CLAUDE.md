# PoolVibes

## Project Overview

Pool maintenance management app. Go CLI with DDD architecture and Datastar hypermedia frontend. Module path: `github.com/joshthewhite/poolvibes`. GitHub repo: `https://github.com/joshthewhite/poolvibes`.

### Structure

- `main.go` — entrypoint, embeds migrations, calls `cmd.Execute()`
- `cmd/` — Cobra CLI commands
  - `cmd/root.go` — root command, Viper config initialization
  - `cmd/serve.go` — serve command, wires all layers (repos, services, server)
- `internal/domain/` — entities, value objects, repository interfaces
- `internal/application/` — command structs, services
- `internal/infrastructure/db/sqlite/` — SQLite repos, connection, migrations
- `internal/interface/web/` — HTTP server, handlers, templates

### CLI & Config

- **CLI framework**: [Cobra](https://github.com/spf13/cobra) — add new commands in `cmd/`
- **Config framework**: [Viper](https://github.com/spf13/viper)
  - Config file: `$HOME/.poolvibes.yaml` or `./.poolvibes.yaml`
  - Override with `--config <path>`
  - Env vars are automatically bound via `viper.AutomaticEnv()`

### Key Dependencies

- **Router**: Go 1.22+ `http.ServeMux` (method-based routing, `r.PathValue()`)
- **Frontend**: [Datastar](https://data-star.dev) — SSE-driven reactive UI via `datastar-go` SDK
- **Templates**: [templ](https://templ.guide) — type-safe HTML templates compiled to Go
- **CSS**: [Bulma](https://bulma.io) v1.0.4 from CDN
- **Database**: SQLite via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO)
- **Migrations**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate) with embedded SQL files

## Development

- **Go version**: 1.25.7
- **Module file**: `go.mod`

### Commands

Uses [Task](https://taskfile.dev) for build automation. Run `task --list` to see all available tasks.

- **Generate templates**: `task templ` (required after editing `.templ` files; generated `*_templ.go` files are committed)
- **Build**: `task build` (outputs to `bin/poolvibes`)
- **Test**: `task test`
- **Test (verbose)**: `task test:verbose`
- **Test (single)**: `go test -v -run TestName ./path/to/package`
- **Lint**: `task lint`
- **Format**: `task fmt`
- **Tidy deps**: `task tidy`
- **Dev server (live reload)**: `task dev` (runs templ generate + air)
- **Build and run**: `task run`
- **Clean**: `task clean`

## Documentation

- **Docs site**: `docs/` directory, built with [Zensical](https://zensical.org), configured in `zensical.toml`
- **README**: `README.md`
- When making code changes, keep `docs/` and `README.md` up to date (CLI flags, features, architecture, dev commands)

## Code Style

- Follow standard Go conventions (effective Go, Go Code Review Comments)
- Use `gofmt` formatting (tabs, standard layout)
- Error handling: return errors, don't panic. Wrap errors with `fmt.Errorf("context: %w", err)`
- Naming: camelCase locals, PascalCase exports, short receiver names
- Keep packages focused and small
- Write table-driven tests
- Use `context.Context` as first parameter where appropriate
