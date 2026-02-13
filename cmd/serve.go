package cmd

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/infrastructure/db/postgres"
	"github.com/joshthewhite/poolvibes/internal/infrastructure/db/sqlite"
	"github.com/joshthewhite/poolvibes/internal/interface/web"
	"github.com/spf13/cobra"
)

var migrationsFS fs.FS

func SetMigrationsFS(f fs.FS) {
	migrationsFS = f
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the PoolVibes web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		dbDSN, _ := cmd.Flags().GetString("db")
		dbDriver, _ := cmd.Flags().GetString("db-driver")

		var (
			db          *sql.DB
			err         error
			chemLogRepo repositories.ChemistryLogRepository
			taskRepo    repositories.TaskRepository
			equipRepo   repositories.EquipmentRepository
			srRepo      repositories.ServiceRecordRepository
			chemRepo    repositories.ChemicalRepository
			userRepo    repositories.UserRepository
			sessionRepo repositories.SessionRepository
		)

		switch dbDriver {
		case "sqlite":
			db, err = sqlite.Open(dbDSN)
			if err != nil {
				return fmt.Errorf("opening database: %w", err)
			}
			defer db.Close()

			if err := sqlite.RunMigrations(db, migrationsFS); err != nil {
				return fmt.Errorf("running migrations: %w", err)
			}

			chemLogRepo = sqlite.NewChemistryLogRepo(db)
			taskRepo = sqlite.NewTaskRepo(db)
			equipRepo = sqlite.NewEquipmentRepo(db)
			srRepo = sqlite.NewServiceRecordRepo(db)
			chemRepo = sqlite.NewChemicalRepo(db)
			userRepo = sqlite.NewUserRepo(db)
			sessionRepo = sqlite.NewSessionRepo(db)

		case "postgres":
			db, err = postgres.Open(dbDSN)
			if err != nil {
				return fmt.Errorf("opening database: %w", err)
			}
			defer db.Close()

			if err := postgres.RunMigrations(db, migrationsFS); err != nil {
				return fmt.Errorf("running migrations: %w", err)
			}

			chemLogRepo = postgres.NewChemistryLogRepo(db)
			taskRepo = postgres.NewTaskRepo(db)
			equipRepo = postgres.NewEquipmentRepo(db)
			srRepo = postgres.NewServiceRecordRepo(db)
			chemRepo = postgres.NewChemicalRepo(db)
			userRepo = postgres.NewUserRepo(db)
			sessionRepo = postgres.NewSessionRepo(db)

		default:
			return fmt.Errorf("unsupported database driver: %s (use 'sqlite' or 'postgres')", dbDriver)
		}

		authSvc := services.NewAuthService(userRepo, sessionRepo)
		userSvc := services.NewUserService(userRepo, sessionRepo)
		chemSvc := services.NewChemistryService(chemLogRepo)
		taskSvc := services.NewTaskService(taskRepo)
		equipSvc := services.NewEquipmentService(equipRepo, srRepo)
		chemicSvc := services.NewChemicalService(chemRepo)

		server := web.NewServer(authSvc, userSvc, chemSvc, taskSvc, equipSvc, chemicSvc)
		return server.Start(addr)
	},
}

func init() {
	serveCmd.Flags().String("addr", ":8080", "server listen address")
	serveCmd.Flags().String("db", defaultDBPath(), "database connection string")
	serveCmd.Flags().String("db-driver", "sqlite", "database driver (sqlite or postgres)")
	rootCmd.AddCommand(serveCmd)
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "poolio.db"
	}
	return home + "/.poolvibes.db"
}
