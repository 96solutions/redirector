package cmd

import (
	"fmt"
	"os"

	"github.com/lroman242/redirector/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "redirector",
	Short: "Redirector is a simple application to handle HTTP redirects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "empty by default. means read from environment")

	rootCmd.PersistentFlags().String("log_level", "info", "set level of logs that should be written")
	rootCmd.PersistentFlags().Bool("log_is_json", true, "should logs be handled in JSON format?")
	rootCmd.PersistentFlags().Bool("log_add_source", true, "should logs contain the source code file/line where a log was created?")
	rootCmd.PersistentFlags().Bool("log_replace_default", true, "should newly created logger replace the default logger?")

	rootCmd.PersistentFlags().String("db_host", "localhost", "storage host")
	rootCmd.PersistentFlags().String("db_port", "3306", "storage port")
	rootCmd.PersistentFlags().String("db_username", "root", "storage username")
	rootCmd.PersistentFlags().String("db_password", "secret", "storage password")
	rootCmd.PersistentFlags().String("db_database", "redirector", "storage database name")

	rootCmd.PersistentFlags().Int("db_conn_max_life", 500, "storage connection max life")
	rootCmd.PersistentFlags().Int("db_max_idle_conn", 50, "storage max idle connections")
	rootCmd.PersistentFlags().Int("db_max_open_conn", 50, "storage max open connections")

	rootCmd.PersistentFlags().String("http_server_host", "", "http server host")
	rootCmd.PersistentFlags().String("http_server_port", "8080", "http server post")
	rootCmd.PersistentFlags().Bool("http_server_ssl", false, "is ssl enabled")
	rootCmd.PersistentFlags().String("http_server_cert", "path/cert.pem", "path to ssl certs")

	rootCmd.PersistentFlags().String("geoip2_db_path", "GeoIP2-City.mmdb", "path to GeoIP2 DB file")

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		panic("Error: " + err.Error())
	}
}

func initConfig() {
	config.GetConfig()
}
