package cmd

import (
	"github.com/visheyra/demo-observability/server"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "appli",
	Long:  "Application",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(appPort, jaegerurl)
	},
}
