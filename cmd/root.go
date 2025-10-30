package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:  "auth-service",
	Long: "Customers service for 'beyachad-maan' web application.",
}

func init() {
	rootCmd.AddCommand(&runCmd)
	rootCmd.AddCommand(&migrateCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
