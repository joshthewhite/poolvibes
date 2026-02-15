# Configuration

PoolVibes can be configured through CLI flags, a config file, or environment variables.

## CLI Flags

### `serve` Command

| Flag | Default | Description |
|------|---------|-------------|
| `--addr` | `:8080` | Server listen address |
| `--db` | `~/.poolvibes.db` | Database connection string |
| `--db-driver` | `sqlite` | Database driver (`sqlite` or `postgres`) |
| `--notify-check-interval` | `1h` | How often to check for due task notifications |
| `--demo` | `false` | Enable demo mode (new non-admin signups get seeded data, auto-expire in 24h) |
| `--demo-max-users` | `50` | Maximum number of concurrent demo users (0 = unlimited) |

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

### PORT

On platforms like Railway that inject a `PORT` environment variable, PoolVibes will automatically use it as the listen port when `--addr` is not explicitly provided. This takes precedence over the default `:8080`.

## Notifications

PoolVibes can send email and SMS notifications when tasks are due. Notifications are checked on a configurable interval (default: 1 hour) and sent at most once per task per day per channel.

### Email (Resend)

| Config Key | Env Var | Description |
|------------|---------|-------------|
| `resend_api_key` | `RESEND_API_KEY` | Resend API key |
| `resend_from` | `RESEND_FROM` | Sender email address (default: `notifications@poolvibes.app`) |

### SMS (Twilio)

| Config Key | Env Var | Description |
|------------|---------|-------------|
| `twilio_account_sid` | `TWILIO_ACCOUNT_SID` | Twilio account SID |
| `twilio_auth_token` | `TWILIO_AUTH_TOKEN` | Twilio auth token |
| `twilio_from_number` | `TWILIO_FROM_NUMBER` | Twilio sender phone number |

### Example Config

```yaml
resend_api_key: "re_..."
resend_from: "notifications@yourdomain.com"
twilio_account_sid: "AC..."
twilio_auth_token: "..."
twilio_from_number: "+15551234567"
notify_check_interval: "1h"
```

Or via environment variables:

```sh
export RESEND_API_KEY="re_..."
export RESEND_FROM="notifications@yourdomain.com"
export TWILIO_ACCOUNT_SID="AC..."
export TWILIO_AUTH_TOKEN="..."
export TWILIO_FROM_NUMBER="+15551234567"
```

Notifications are only enabled when the corresponding API keys are configured. Users can toggle email/SMS preferences and set their phone number from the Settings tab in the app.

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
