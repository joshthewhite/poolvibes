# Configuration

PoolVibes can be configured through CLI flags, a config file, or environment variables.

## CLI Flags

### `serve` Command

| Flag | Default | Description |
|------|---------|-------------|
| `--addr` | `:8080` | Server listen address |
| `--db` | `~/.poolvibes.db` | Database connection string |
| `--db-driver` | `sqlite` | Database driver (`sqlite` or `postgres`) |

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | (see below) | Path to config file |

## Config File

PoolVibes uses [Viper](https://github.com/spf13/viper) for configuration. It searches for a `.poolvibes.yaml` file in:

1. `$HOME/.poolvibes.yaml`
2. `./.poolvibes.yaml` (current working directory)

You can also specify a config file explicitly:

```sh
./poolvibes serve --config /path/to/config.yaml
```

### Example Config (SQLite)

```yaml
addr: ":3000"
db: "/var/lib/poolvibes/pool.db"
```

### Example Config (PostgreSQL)

```yaml
addr: ":3000"
db-driver: "postgres"
db: "postgres://user:pass@localhost:5432/poolvibes?sslmode=disable"
```

## Environment Variables

All configuration options can be set via environment variables. Viper's `AutomaticEnv()` binds them automatically:

```sh
export ADDR=":3000"
export DB="/var/lib/poolvibes/pool.db"
./poolvibes serve
```

## Database

PoolVibes supports two database backends: **SQLite** (default) and **PostgreSQL**.

### Migrations

Database migrations run automatically on server startup. Migration files are embedded in the binary at build time using Go's `embed` package, so no external migration files are needed. Separate migration sets are maintained for each database driver.

### SQLite (default)

SQLite is the default database. The database file is created automatically on first run.

The following SQLite pragmas are configured automatically:

- **WAL mode** (`journal_mode=WAL`) — Enables concurrent reads during writes
- **Foreign keys** (`foreign_keys=1`) — Enforces referential integrity

### PostgreSQL

For hosted or multi-instance deployments, PostgreSQL can be used instead:

```sh
./poolvibes serve --db-driver postgres --db "postgres://user:pass@localhost:5432/poolvibes?sslmode=disable"
```

PostgreSQL uses native types (`UUID`, `TIMESTAMPTZ`, `BOOLEAN`, `DOUBLE PRECISION`) for better type safety and performance. The PostgreSQL database must exist before starting the server; tables and migrations are applied automatically.
