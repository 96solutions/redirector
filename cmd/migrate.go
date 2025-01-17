package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/lroman242/redirector/config"
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
		db, err := sql.Open("mysql", cfg.DBConf.DSN()+"?multiStatements=true")
		if err != nil {
			panic(fmt.Errorf("unable to connect to the database. error: %w", err))
		}

		err = db.Ping()
		if err != nil {
			panic(fmt.Errorf("unsuccessfull ping database. error: %w", err))
		}

		driver, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			panic(err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			"file://./migrations",
			"mysql",
			driver,
		)
		if err != nil {
			panic(err)
		}

		if steps != 0 {
			panic(m.Steps(steps))
		} else {
			panic(m.Steps(steps))
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	rootCmd.PersistentFlags().IntVar(&steps, "steps", 1, "number of steps to migrate or rollback")
}
