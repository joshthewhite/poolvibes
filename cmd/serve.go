package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/infrastructure/db/postgres"
	"github.com/joshthewhite/poolvibes/internal/infrastructure/db/sqlite"
	"github.com/joshthewhite/poolvibes/internal/infrastructure/notify"
	"github.com/joshthewhite/poolvibes/internal/interface/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrationsFS fs.FS

func SetMigrationsFS(f fs.FS) {
	migrationsFS = f
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the PoolVibes web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := viper.GetString("addr")
		if port := os.Getenv("PORT"); port != "" && !cmd.Flags().Changed("addr") {
			addr = ":" + port
		}
		dbDSN := viper.GetString("db")
		dbDriver := viper.GetString("db-driver")

		var (
			db            *sql.DB
			err           error
			chemLogRepo   repositories.ChemistryLogRepository
			taskRepo      repositories.TaskRepository
			equipRepo     repositories.EquipmentRepository
			srRepo        repositories.ServiceRecordRepository
			chemRepo      repositories.ChemicalRepository
			userRepo      repositories.UserRepository
			sessionRepo   repositories.SessionRepository
			taskNotifRepo repositories.TaskNotificationRepository
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
			taskNotifRepo = sqlite.NewTaskNotificationRepo(db)

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
			taskNotifRepo = postgres.NewTaskNotificationRepo(db)

		default:
			return fmt.Errorf("unsupported database driver: %s (use 'sqlite' or 'postgres')", dbDriver)
		}

		authSvc := services.NewAuthService(userRepo, sessionRepo)
		userSvc := services.NewUserService(userRepo, sessionRepo)
		chemSvc := services.NewChemistryService(chemLogRepo)
		taskSvc := services.NewTaskService(taskRepo)
		equipSvc := services.NewEquipmentService(equipRepo, srRepo)
		chemicSvc := services.NewChemicalService(chemRepo)

		// Set up notification service
		var emailNotifier services.Notifier
		var smsNotifier services.Notifier

		if apiKey := viper.GetString("resend_api_key"); apiKey != "" {
			from := viper.GetString("resend_from")
			if from == "" {
				from = "notifications@poolvibes.app"
			}
			emailNotifier = notify.NewResendNotifier(apiKey, from)
			log.Println("Email notifications enabled (Resend)")
		}

		if sid := viper.GetString("twilio_account_sid"); sid != "" {
			token := viper.GetString("twilio_auth_token")
			fromNum := viper.GetString("twilio_from_number")
			if token != "" && fromNum != "" {
				smsNotifier = notify.NewTwilioNotifier(sid, token, fromNum)
				log.Println("SMS notifications enabled (Twilio)")
			}
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		if emailNotifier != nil || smsNotifier != nil {
			intervalStr := viper.GetString("notify-check-interval")
			interval, err := time.ParseDuration(intervalStr)
			if err != nil {
				interval = 1 * time.Hour
			}
			notifSvc := services.NewNotificationService(taskRepo, userRepo, taskNotifRepo, emailNotifier, smsNotifier, interval)
			go notifSvc.Start(ctx)
		}

		server := web.NewServer(authSvc, userSvc, chemSvc, taskSvc, equipSvc, chemicSvc)
		return server.Start(ctx, addr)
	},
}

func init() {
	serveCmd.Flags().String("addr", ":8080", "server listen address")
	serveCmd.Flags().String("db", defaultDBPath(), "database connection string")
	serveCmd.Flags().String("db-driver", "sqlite", "database driver (sqlite or postgres)")
	serveCmd.Flags().String("notify-check-interval", "1h", "how often to check for due task notifications")

	viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
	viper.BindPFlag("db", serveCmd.Flags().Lookup("db"))
	viper.BindPFlag("db-driver", serveCmd.Flags().Lookup("db-driver"))
	viper.BindPFlag("notify-check-interval", serveCmd.Flags().Lookup("notify-check-interval"))

	rootCmd.AddCommand(serveCmd)
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "poolio.db"
	}
	return home + "/.poolvibes.db"
}
