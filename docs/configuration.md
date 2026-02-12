# Configuration

PoolVibes can be configured through CLI flags, a config file, or environment variables.

## CLI Flags

### `serve` Command

| Flag | Default | Description |
|------|---------|-------------|
| `--addr` | `:8080` | Server listen address |
| `--db` | `~/.poolio.db` | SQLite database file path |

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
./poolio serve --config /path/to/config.yaml
```

### Example Config

```yaml
addr: ":3000"
db: "/var/lib/poolvibes/pool.db"
```

## Environment Variables

All configuration options can be set via environment variables. Viper's `AutomaticEnv()` binds them automatically:

```sh
export ADDR=":3000"
export DB="/var/lib/poolvibes/pool.db"
./poolio serve
```

## Database

PoolVibes uses SQLite as its database. The database file is created automatically on first run.

### Migrations

Database migrations run automatically on server startup. Migration files are embedded in the binary at build time using Go's `embed` package, so no external migration files are needed.

### SQLite Settings

The following SQLite pragmas are configured automatically:

- **WAL mode** (`journal_mode=WAL`) — Enables concurrent reads during writes
- **Foreign keys** (`foreign_keys=1`) — Enforces referential integrity
