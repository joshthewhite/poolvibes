# PoolVibes

A pool maintenance management app built with Go, following Domain-Driven Design patterns with a [Datastar](https://data-star.dev) hypermedia frontend.

## Features

- **Dashboard** — At-a-glance overview with water quality summary, task status, low stock alerts, and pH/chlorine trend charts (Chart.js).
- **Authentication** — Email/password sign-up and sign-in with cookie-based sessions. Per-user data isolation (multi-tenancy).
- **Admin Panel** — Admin users can manage accounts (enable/disable users, grant admin access).
- **Water Chemistry** — Log pH, free/combined chlorine, total alkalinity, CYA, calcium hardness, and temperature. Out-of-range values are highlighted automatically. Server-side pagination with sortable columns and date/out-of-range filters. Generate treatment plans with chemical dosages based on your pool size.
- **Task Scheduling** — Create recurring maintenance tasks (daily, weekly, monthly). Completing a task auto-generates the next occurrence.
- **Equipment Tracking** — Track pool equipment with categories, manufacturer info, warranty status, and service history.
- **Chemical Inventory** — Monitor chemical stock levels with low-stock alerts and quick-adjust buttons.
- **Notifications** — Email (Resend) and SMS (Twilio) alerts when tasks are due. Per-user preferences via Settings tab.
- **Demo Mode** — Enable `--demo` to let potential customers sign up and see the app pre-populated with a year of realistic data. Demo users auto-expire after 24 hours. Admins can convert demo users to regular accounts.

## Tech Stack

- **Go** with [Cobra](https://github.com/spf13/cobra) CLI and Go 1.22+ `http.ServeMux` router
- **Datastar** for reactive SSE-driven UI (no JavaScript framework)
- **[templ](https://templ.guide)** for type-safe HTML templates compiled to Go
- **SQLite** via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO) — default
- **PostgreSQL** via [pgx](https://github.com/jackc/pgx) — optional, for hosted deployments
- **Bulma** CSS from CDN
- **Resend** for email notifications, **Twilio** for SMS notifications
- **DDD architecture** — domain entities, repository interfaces, application services, infrastructure implementations

## Getting Started

```sh
go build -o poolvibes .
./poolvibes serve
```

Open http://localhost:8080 — you'll be redirected to sign up on first visit.

### Options

```
--addr string                  server listen address (default ":8080")
--db string                    database connection string (default "~/.poolvibes.db")
--db-driver string             database driver: sqlite or postgres (default "sqlite")
--notify-check-interval string how often to check for due task notifications (default "1h")
--demo                         enable demo mode (default false)
--demo-max-users int           max concurrent demo users (default 50, 0 = unlimited)
```

Database migrations run automatically on startup.

#### PostgreSQL

To use PostgreSQL instead of SQLite:

```sh
./poolvibes serve --db-driver postgres --db "postgres://user:pass@localhost:5432/poolvibes?sslmode=disable"
```

### Notifications

To enable task due notifications, configure your API keys in `~/.poolvibes.yaml`:

```yaml
resend_api_key: "re_..."
resend_from: "notifications@yourdomain.com"
twilio_account_sid: "AC..."
twilio_auth_token: "..."
twilio_from_number: "+15551234567"
```

Or via environment variables (`RESEND_API_KEY`, `TWILIO_ACCOUNT_SID`, etc.). Notifications are only sent when the corresponding keys are configured. Users can toggle email/SMS preferences from the Settings tab.

## Deployment (Railway)

PoolVibes can be deployed to [Railway](https://railway.com) with PostgreSQL:

1. Create a new project on Railway and add a **PostgreSQL** service
2. Add a service connected to your GitHub repo — Railway auto-detects the `Dockerfile`
3. Set these service variables:
   - `DB_DRIVER` = `postgres`
   - `DB` = `${{Postgres.DATABASE_URL}}`
4. Enable **Wait for CI** in service settings so Railway waits for GitHub Actions to pass before deploying

The app automatically uses Railway's `PORT` environment variable. Migrations run on startup.

## Project Structure

```
poolvibes/
├── main.go                          # entrypoint, embeds migrations
├── cmd/
│   ├── root.go                      # Cobra root command
│   └── serve.go                     # serve command, wires all layers
├── migrations/
│   ├── sqlite/                      # SQLite migrations (embedded)
│   └── postgres/                    # PostgreSQL migrations (embedded)
└── internal/
    ├── domain/
    │   ├── entities/                # User, Session, ChemistryLog, Task, Equipment, etc.
    │   ├── valueobjects/            # Recurrence, Quantity
    │   └── repositories/            # interfaces
    ├── application/
    │   ├── command/                 # CRUD command structs
    │   └── services/                # business logic
    ├── infrastructure/
    │   ├── db/
    │   │   ├── sqlite/              # SQLite repos + connection
    │   │   └── postgres/            # PostgreSQL repos + connection
    │   └── notify/                  # Email (Resend) and SMS (Twilio) notifiers
    └── interface/
        └── web/
            ├── server.go            # HTTP server + routes
            ├── handlers/            # SSE handlers per feature
            └── templates/           # templ components (*.templ + generated *_templ.go)
```

## Development

```sh
templ generate          # regenerate templates (after editing .templ files)
go build ./...          # build
go vet ./...            # lint
go test ./...           # test
gofmt -w .              # format
go mod tidy             # tidy deps
```
