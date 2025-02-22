// Package cmd contains application commands and their implementations.
package cmd

import (
	"os"

	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/infrastructure/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
// It provides the main entry point and common configuration for all subcommands.
var rootCmd = &cobra.Command{
	Use:   "redirector",
	Short: "Redirector is a simple application to handle HTTP redirects.",
	Long: `Redirector is a microservice that handles HTTP redirects based on configurable rules.
It supports various redirect types including geo-targeting, device detection, and A/B testing.`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := config.GetConfig()

		log := logger.NewLogger(cfg.LogConf)
		log.Info("Hello ...")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// If an error occurs, the application will exit with status code 1.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init initializes the command configuration and flags.
// It sets up all configuration flags and binds them to viper for config management.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "empty by default. means read from environment")

	rootCmd.PersistentFlags().String("log_level", "info", "set level of logs that should be written")
	rootCmd.PersistentFlags().Bool("log_is_json", true, "should logs be handled in JSON format?")
	rootCmd.PersistentFlags().Bool(
		"log_add_source",
		true,
		"should logs contain the source code file/line where a log was created?",
	)
	rootCmd.PersistentFlags().Bool(
		"log_replace_default",
		true,
		"should newly created logger replace the default logger?",
	)

	// OpenSearch configuration flags
	rootCmd.PersistentFlags().String("log_open_search_host", "localhost", "OpenSearch node host")
	rootCmd.PersistentFlags().String("log_open_search_port", "9200", "OpenSearch node port")
	rootCmd.PersistentFlags().String("log_open_search_index", "redirector", "OpenSearch index name")
	rootCmd.PersistentFlags().String("log_open_search_user", "logger", "OpenSearch user (auth)")
	rootCmd.PersistentFlags().String("log_open_search_pass", "", "OpenSearch password (auth)")

	// Database configuration flags
	rootCmd.PersistentFlags().String("db_host", "localhost", "storage host")
	rootCmd.PersistentFlags().String("db_port", "3306", "storage port")
	rootCmd.PersistentFlags().String("db_username", "root", "storage username")
	rootCmd.PersistentFlags().String("db_password", "secret", "storage password")
	rootCmd.PersistentFlags().String("db_database", "redirector", "storage database name")

	// SQL Storage configuration flags
	rootCmd.PersistentFlags().Int("db_conn_max_life", 500, "storage connection max life")
	rootCmd.PersistentFlags().Int("db_max_idle_conn", 50, "storage max idle connections")
	rootCmd.PersistentFlags().Int("db_max_open_conn", 50, "storage max open connections")

	// HTTP server configuration flags
	rootCmd.PersistentFlags().String("http_server_host", "", "http server host")
	rootCmd.PersistentFlags().String("http_server_port", "8080", "http server post")
	rootCmd.PersistentFlags().Bool("http_server_ssl", false, "is ssl enabled")
	rootCmd.PersistentFlags().String("http_server_cert", "path/cert.pem", "path to ssl certs")

	// Redis configuration flags
	rootCmd.PersistentFlags().String("redis_host", "localhost", "Redis server hostname")
	rootCmd.PersistentFlags().String("redis_port", "6379", "Redis server port")
	rootCmd.PersistentFlags().String("redis_pass", "", "Redis server password")
	rootCmd.PersistentFlags().Int("redis_db", 0, "Redis database number")
	rootCmd.PersistentFlags().Int("redis_max_retries", 3, "Maximum number of retries before giving up")
	rootCmd.PersistentFlags().Int("redis_min_retry_backoff", 8, "Minimum backoff between each retry in milliseconds")
	rootCmd.PersistentFlags().Int("redis_max_retry_backoff", 512, "Maximum backoff between each retry in milliseconds")
	rootCmd.PersistentFlags().Int("redis_dial_timeout", 5, "Timeout for establishing new connections in seconds")
	rootCmd.PersistentFlags().Int("redis_read_timeout", 3, "Timeout for socket reads in seconds")
	rootCmd.PersistentFlags().Int("redis_write_timeout", 3, "Timeout for socket writes in seconds")
	rootCmd.PersistentFlags().Int("redis_pool_size", 10, "Maximum number of socket connections")
	rootCmd.PersistentFlags().Int(
		"redis_pool_timeout",
		4,
		"Time client waits for connection if all connections are busy in seconds",
	)

	// GeoIP2 configuration flags
	rootCmd.PersistentFlags().String("geoip2_db_path", "GeoIP2-City.mmdb", "path to GeoIP2 DB file")

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		panic("Error: " + err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
// It's called before executing any command to ensure proper configuration.
func initConfig() {
	config.GetConfig()
}
