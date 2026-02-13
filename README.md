# PoolVibes

A pool maintenance management app built with Go, following Domain-Driven Design patterns with a [Datastar](https://data-star.dev) hypermedia frontend.

## Features

- **Authentication** — Email/password sign-up and sign-in with cookie-based sessions. Per-user data isolation (multi-tenancy).
- **Admin Panel** — Admin users can manage accounts (enable/disable users, grant admin access).
- **Water Chemistry** — Log pH, free/combined chlorine, total alkalinity, CYA, calcium hardness, and temperature. Out-of-range values are highlighted automatically.
- **Task Scheduling** — Create recurring maintenance tasks (daily, weekly, monthly). Completing a task auto-generates the next occurrence.
- **Equipment Tracking** — Track pool equipment with categories, manufacturer info, warranty status, and service history.
- **Chemical Inventory** — Monitor chemical stock levels with low-stock alerts and quick-adjust buttons.

## Tech Stack

- **Go** with [Cobra](https://github.com/spf13/cobra) CLI and Go 1.22+ `http.ServeMux` router
- **Datastar** for reactive SSE-driven UI (no JavaScript framework)
- **SQLite** via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO)
- **Bulma** CSS from CDN
- **DDD architecture** — domain entities, repository interfaces, application services, infrastructure implementations

## Getting Started

```sh
go build -o poolvibes .
./poolvibes serve
```

Open http://localhost:8080 — you'll be redirected to sign up on first visit.

### Options

```
--addr string   server listen address (default ":8080")
--db string     SQLite database path (default "~/.poolvibes.db")
```

Database migrations run automatically on startup.

## Project Structure

```
poolvibes/
├── main.go                          # entrypoint, embeds migrations
├── cmd/
│   ├── root.go                      # Cobra root command
│   └── serve.go                     # serve command, wires all layers
├── migrations/                      # SQLite migrations (embedded)
└── internal/
    ├── domain/
    │   ├── entities/                # User, Session, ChemistryLog, Task, Equipment, etc.
    │   ├── valueobjects/            # Recurrence, Quantity
    │   └── repositories/            # interfaces
    ├── application/
    │   ├── command/                 # CRUD command structs
    │   └── services/                # business logic
    ├── infrastructure/
    │   └── db/sqlite/               # SQLite repos + migrations
    └── interface/
        └── web/
            ├── server.go            # HTTP server + routes
            ├── handlers/            # SSE handlers per feature
            └── templates/           # layout.html
```

## Development

```sh
go build ./...          # build
go vet ./...            # lint
go test ./...           # test
gofmt -w .              # format
go mod tidy             # tidy deps
```
