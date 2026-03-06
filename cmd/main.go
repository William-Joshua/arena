package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version information can be injected at build time via ldflags:
//
//	go build -ldflags="-X main.version=1.0.0 -X main.commit=abc -X main.date=2024-01-01"
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "arena",
		Short: "Arena application",
	}

	rootCmd.AddCommand(newApiserverCmd())
	rootCmd.AddCommand(newMigrateCmd())
	rootCmd.AddCommand(newVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
