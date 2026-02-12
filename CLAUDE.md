# Poolio

## Project Overview

Go CLI application. Module path: `github.com/josh/poolio`.

### Structure

- `main.go` — entrypoint, calls `cmd.Execute()`
- `cmd/` — Cobra CLI commands (add new subcommands here)
  - `cmd/root.go` — root command, Viper config initialization

### CLI & Config

- **CLI framework**: [Cobra](https://github.com/spf13/cobra) — add new commands in `cmd/`
- **Config framework**: [Viper](https://github.com/spf13/viper)
  - Config file: `$HOME/.poolio.yaml` or `./.poolio.yaml`
  - Override with `--config <path>`
  - Env vars are automatically bound via `viper.AutomaticEnv()`

## Development

- **Go version**: 1.23.0
- **Module file**: `go.mod`

### Commands

- **Build**: `go build ./...`
- **Test**: `go test ./...`
- **Test (verbose)**: `go test -v ./...`
- **Test (single)**: `go test -v -run TestName ./path/to/package`
- **Lint**: `go vet ./...`
- **Format**: `gofmt -w .`
- **Tidy deps**: `go mod tidy`

## Code Style

- Follow standard Go conventions (effective Go, Go Code Review Comments)
- Use `gofmt` formatting (tabs, standard layout)
- Error handling: return errors, don't panic. Wrap errors with `fmt.Errorf("context: %w", err)`
- Naming: camelCase locals, PascalCase exports, short receiver names
- Keep packages focused and small
- Write table-driven tests
- Use `context.Context` as first parameter where appropriate
