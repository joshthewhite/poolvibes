# Taskfile Design

**Date:** 2026-02-15
**Status:** Approved

## Summary

Add a [Taskfile](https://taskfile.dev) (v3) to standardize build, test, lint, format, and dev workflows. Introduce a `bin/` output directory for built artifacts and [air](https://github.com/air-verse/air) for live reload during development.

## Build Output

- All binaries output to `bin/` (e.g. `bin/poolvibes`)
- `.gitignore` updated: replace individual binary names (`poolio`, `poolvibes`) with `bin/`
- Dockerfile updated to build to `bin/poolvibes`

## Tasks

| Task | Command | Description |
|------|---------|-------------|
| `build` | `go build -o bin/poolvibes .` | Build the binary |
| `test` | `go test ./...` | Run all tests |
| `test:verbose` | `go test -v ./...` | Run tests with verbose output |
| `lint` | `go vet ./...` | Run linter |
| `fmt` | `gofmt -w .` | Format code |
| `templ` | `templ generate` | Generate templ templates |
| `dev` | templ generate, then air | Full dev loop with live reload |
| `run` | Build then execute `bin/poolvibes serve` | Quick build-and-run |
| `clean` | Remove `bin/` | Clean build artifacts |
| `docker:build` | `docker build .` | Build Docker image |
| `docker:up` | `docker-compose up` | Start all services |
| `docker:down` | `docker-compose down` | Stop all services |
| `tidy` | `go mod tidy` | Tidy dependencies |

## Air Configuration

- `.air.toml` at project root
- Watch `.go` and `.templ` files
- Run `templ generate` before rebuild
- Build output to `bin/poolvibes`
- Exclude `_templ.go`, `_test.go`, vendor, bin, e2e, docs directories

## Other Changes

- Update `.gitignore`: replace `poolio`/`poolvibes` with `bin/`
- Update `CLAUDE.md` and `README.md` to reference `task` commands
- Update `docs/development.md` with Taskfile usage
- CI workflow keeps raw Go commands (no Task dependency in CI)
