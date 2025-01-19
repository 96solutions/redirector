package cmd

import (
	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/infrastructure/logger"
	"github.com/spf13/cobra"
)

// configCmd represents the config command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show application configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		log := logger.NewLogger(cfg.LogConf)
		log.Info("Application config", "config", cfg)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
