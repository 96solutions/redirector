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

	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().String("config", ".env", "config file (default is $HOME/.env)")

	rootCmd.PersistentFlags().String("log_level", "warning", "set level of logs that should be written")
	rootCmd.PersistentFlags().String("log_dir", ".", "set path to the directory where logs should be written")
	rootCmd.PersistentFlags().String("log_file", "redirector.log", "set file name where logs should be written")

	rootCmd.PersistentFlags().String("mysql_host", "localhost", "storage host")
	rootCmd.PersistentFlags().String("mysql_port", "3306", "storage port")
	rootCmd.PersistentFlags().String("mysql_username", "root", "storage username")
	rootCmd.PersistentFlags().String("mysql_password", "secret", "storage password")
	rootCmd.PersistentFlags().String("mysql_database", "redirector", "storage database name")

	rootCmd.PersistentFlags().String("http_server_host", "localhost", "http server host")
	rootCmd.PersistentFlags().String("http_server_port", "8080", "http server post")
	rootCmd.PersistentFlags().Bool("http_server_ssl", false, "is ssl enabled")
	rootCmd.PersistentFlags().String("http_server_cert", "path/cert.pem", "path to ssl certs")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	config.GetConfig()
}
