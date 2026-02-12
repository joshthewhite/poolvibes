package main

import (
	"embed"

	"github.com/josh/poolio/cmd"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	cmd.SetMigrationsFS(migrationsFS)
	cmd.Execute()
}
