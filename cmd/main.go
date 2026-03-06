package main

import (
	"github.com/spf13/cobra"
	"os"
)

// Setup the Cobra root command
func main() {
	var rootCmd = &cobra.Command{
		Use: "app",
		Short: "An application that does something",
	}

	// Adding subcommands
	rootCmd.AddCommand(apiserverCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(versionCmd)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Placeholders for actual commands
var apiserverCmd = &cobra.Command{
	Use: "apiserver",
	Short: "Run the API server",
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Short: "Run the database migrations",
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version",
}