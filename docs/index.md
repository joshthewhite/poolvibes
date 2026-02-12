# PoolVibes

Pool maintenance management app built with Go, following Domain-Driven Design patterns with a [Datastar](https://data-star.dev) hypermedia frontend.

## Features

- **[Water Chemistry](features/water-chemistry.md)** — Log pH, chlorine, alkalinity, CYA, calcium hardness, and temperature with automatic out-of-range highlighting.
- **[Task Scheduling](features/tasks.md)** — Create recurring maintenance tasks that auto-generate the next occurrence on completion.
- **[Equipment Tracking](features/equipment.md)** — Track pool equipment with warranty status and service history.
- **[Chemical Inventory](features/chemicals.md)** — Monitor chemical stock levels with low-stock alerts and quick-adjust buttons.

## Quick Start

```sh
go build -o poolvibes .
./poolvibes serve
```

Open [http://localhost:8080](http://localhost:8080) to access the app.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25+ |
| CLI | [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper) |
| Router | Go `http.ServeMux` (method-based routing) |
| Frontend | [Datastar](https://data-star.dev) (SSE-driven reactive UI) |
| CSS | [Bulma](https://bulma.io) 1.0.4 |
| Database | SQLite via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) with embedded SQL |
| Architecture | Domain-Driven Design |
