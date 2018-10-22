package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{Use: "app", Short: "appli", Long: "Application"}

var jaegerurl string
var appPort int

func init() {
	RootCmd.PersistentFlags().StringVar(&jaegerurl, "jaegerurl", "", "Set jaegger collector endpoint")
	RootCmd.PersistentFlags().IntVar(&appPort, "port", 8080, "set application listen port")
	RootCmd.AddCommand(startCmd)
}
