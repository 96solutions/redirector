package cmd

import (
	"log/slog"
	"os"

	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/registry"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Long:  `Start HTTP server`,
	Run: func(_ *cobra.Command, args []string) {
		reg := registry.NewRegistry(config.GetConfig())

		server := reg.NewServer()

		err := server.Start()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
