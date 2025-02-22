// Package cmd contains application commands and their implementations.
package cmd

import (
	"log/slog"
	"strconv"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/registry"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command.
// It handles database schema migrations using migration files.
var migrateCmd = &cobra.Command{
	Use:   "migrate [steps]",
	Short: "Roll forward or roll back database schema according to the migration files.",
	Long: `Execute database migrations to update or rollback the schema.
Positive steps value migrates forward, negative steps rolls back migrations.
Migration files must be present in the ./migrations directory.

Examples:
  # Apply the next migration
  redirector migrate 1

  # Rollback the last migration
  redirector migrate -1

  # Apply the next 3 migrations
  redirector migrate 3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse steps argument
		steps, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("Invalid steps value", "error", err)
			return
		}

		if steps == 0 {
			slog.Info("Steps value is 0, nothing to do.")
			return
		}

		cfg := config.GetConfig()
		r := registry.NewRegistry(cfg)

		db := r.NewDB()
		defer db.Close()

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			slog.Error("Failed to create postgres driver", "error", err)
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://./migrations",
			"postgres",
			driver,
		)
		if err != nil {
			slog.Error("Failed to create migration instance", "error", err)
			return
		}

		// Apply migrations
		if err := m.Steps(steps); err != nil {
			if err == migrate.ErrNoChange {
				slog.Info("No migrations to apply")
				return
			}
			slog.Error("Failed to apply migrations", "error", err)
			return
		}

		// Get current version
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			slog.Error("Failed to get migration version", "error", err)
			return
		}

		if dirty {
			slog.Warn("Database is in dirty state", "version", version)
			return
		}

		slog.Info("Migrations applied successfully",
			"direction", getDirection(steps),
			"steps", abs(steps),
			"version", version,
		)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

// getDirection returns a string indicating the migration direction.
func getDirection(steps int) string {
	if steps > 0 {
		return "up"
	}
	return "down"
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
