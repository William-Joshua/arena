package main

import (
	"fmt"

	"github.com/spf13/cobra"

	internalmigrate "cc.io/arena/internal/migrate"
)

func newMigrateCmd() *cobra.Command {
	var (
		databaseURL   string
		migrationsDir string
	)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if databaseURL == "" {
				return fmt.Errorf("--database-url is required")
			}
			return internalmigrate.Run(internalmigrate.Config{
				DatabaseURL:   databaseURL,
				MigrationsDir: migrationsDir,
			})
		},
	}

	cmd.Flags().StringVar(&databaseURL, "database-url", "", "Postgres connection URL (required)")
	cmd.Flags().StringVar(&migrationsDir, "migrations-dir", "migrations", "path to SQL migration files")
	return cmd
}
