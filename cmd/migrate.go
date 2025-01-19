package cmd

import (
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/registry"
	"github.com/spf13/cobra"
)

var steps int

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Roll forward or roll back database schema according to the migration files",
	Args: func(cmd *cobra.Command, args []string) error {
		//TODO: validate args

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		r := registry.NewRegistry(cfg)

		db := r.NewDB()

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			panic(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://./migrations",
			"postgres",
			driver,
		)
		if err != nil {
			panic(err)
		}

		if steps != 0 {
			err = m.Steps(steps)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	rootCmd.PersistentFlags().IntVar(&steps, "steps", 1, "number of steps to migrate or rollback")
}
