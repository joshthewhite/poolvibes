# Getting Started

## Installation

Clone the repository and build:

```sh
git clone https://github.com/joshthewhite/poolvibes.git
cd poolvibes
go build -o poolvibes .
```

This produces a single `poolvibes` binary with no external dependencies — the SQLite driver is pure Go and migrations are embedded in the binary.

## Running the Server

Start the web server with the default settings:

```sh
./poolvibes serve
```

This will:

1. Create (or open) a SQLite database at `~/.poolvibes.db`
2. Run any pending database migrations automatically
3. Start the HTTP server on port 8080

Open [http://localhost:8080](http://localhost:8080) in your browser.

### Custom Address and Database

```sh
./poolvibes serve --addr :3000 --db ./mypool.db
```

See [Configuration](configuration.md) for all available options.

## First Use

When you first open the app, you'll see an empty dashboard with four tabs:

1. **Water Chemistry** — Start by logging your first water test. Click "New" and enter your readings.
2. **Tasks** — Set up recurring maintenance tasks like "Check chlorine" or "Clean filter."
3. **Equipment** — Add your pool equipment (pump, filter, heater, etc.) with warranty and service info.
4. **Chemicals** — Track your chemical inventory and set low-stock alert thresholds.

The app uses a single-page interface powered by [Datastar](https://data-star.dev). All interactions happen through Server-Sent Events — no page reloads needed.
