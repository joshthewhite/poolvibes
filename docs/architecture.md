# Architecture

PoolVibes follows Domain-Driven Design (DDD) with a layered architecture. Dependencies point inward — outer layers depend on inner layers, never the reverse.

## Layers

```
┌─────────────────────────────────────┐
│           Interface (Web)           │  HTTP handlers, templates
├─────────────────────────────────────┤
│          Infrastructure             │  SQLite repos, DB connection
├─────────────────────────────────────┤
│           Application               │  Commands, services
├─────────────────────────────────────┤
│             Domain                  │  Entities, value objects, repo interfaces
└─────────────────────────────────────┘
```

### Domain

Pure business logic with no external dependencies. Contains:

- **Entities** — `ChemistryLog`, `Task`, `Equipment`, `ServiceRecord`, `Chemical` with validation rules and business methods
- **Value Objects** — `Recurrence` (frequency + interval with next-due-date calculation), `Quantity` (amount + unit)
- **Repository Interfaces** — Abstractions that infrastructure implements

### Application

Orchestrates domain logic through:

- **Commands** — CRUD command structs (DTOs) for each feature
- **Services** — Business logic coordination (e.g., auto-rescheduling tasks on completion, stock adjustment validation)

### Infrastructure

External concerns:

- **SQLite Repositories** — Implement domain repository interfaces using SQL queries
- **Connection** — Database connection management, migration runner
- **Migrations** — SQL files embedded in the binary via Go's `embed` package

### Interface

User-facing layer:

- **HTTP Server** — Go `http.ServeMux` with method-based routing (`GET /tasks`, `POST /tasks`, etc.)
- **Handlers** — SSE handlers using Datastar for reactive UI updates
- **Templates** — Single `layout.html` serving the SPA shell

## Directory Structure

```
poolvibes/
├── main.go                          # Entrypoint, embeds migrations
├── cmd/
│   ├── root.go                      # Cobra root command, Viper config
│   └── serve.go                     # Serve command, wires all layers
├── migrations/                      # SQLite migrations (embedded)
└── internal/
    ├── domain/
    │   ├── entities/                # ChemistryLog, Task, Equipment, ServiceRecord, Chemical
    │   ├── valueobjects/            # Recurrence, Quantity
    │   └── repositories/            # Interfaces
    ├── application/
    │   ├── command/                 # CRUD command structs
    │   └── services/                # Business logic
    ├── infrastructure/
    │   └── db/sqlite/               # SQLite repos, connection, migrations
    └── interface/
        └── web/
            ├── server.go            # HTTP server + routes
            ├── handlers/            # SSE handlers per feature
            └── templates/           # layout.html
```

## Tech Stack

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Router | `http.ServeMux` | Go 1.22+ method routing, no external dependency |
| Frontend | Datastar | SSE-driven reactive UI, no JavaScript framework |
| CSS | Bulma 1.0.4 | Lightweight, CDN-hosted |
| Database | modernc.org/sqlite | Pure Go, no CGO required |
| Migrations | golang-migrate | Embedded SQL files, auto-run on startup |
| CLI | Cobra + Viper | Standard Go CLI pattern |

## Request Flow

1. Browser sends request (or Datastar sends SSE request via `data-*` attributes)
2. `http.ServeMux` routes to the appropriate handler
3. Handler parses the request into a command struct
4. Service validates and executes business logic via repository interfaces
5. SQLite repository performs the database operation
6. Handler sends SSE response back, patching the UI via Datastar
