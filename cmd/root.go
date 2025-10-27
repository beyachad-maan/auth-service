package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:  "anglit-customers-service",
	Long: "Customers service for 'Anglit' mobile application.",
}

func init() {
	rootCmd.AddCommand(&runCmd)
	rootCmd.AddCommand(&migrateCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
