// Package cmd contains application commands and their implementations.
package cmd

import (
	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/infrastructure/logger"
	"github.com/spf13/cobra"
)

// configCmd represents the config command.
// It displays the current application configuration in a human-readable format.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show application configuration.",
	Long: "Display the current application configuration including all settings" +
		" from environment variables and config files.",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := config.GetConfig()

		log := logger.NewLogger(cfg.LogConf)
		log.Info("Application config", "config", cfg)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
