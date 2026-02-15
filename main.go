package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/joshthewhite/poolvibes/cmd"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cmd.SetMigrationsFS(migrationsFS)
	cmd.Execute()
}
