package cmd

import (
	"fmt"

	"github.com/lroman242/redirector/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show application configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := fmt.Printf("Application config:\n%s\n", config.GetConfig())
		if err != nil {
			panic("Error: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
