package main

import (
	"embed"

	"github.com/joshthewhite/poolvibes/cmd"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	cmd.SetMigrationsFS(migrationsFS)
	cmd.Execute()
}
