// Package migrate provides a database migration runner.
package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // pgx v5 driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // file source
)

// Config holds the parameters needed to connect to the database and locate
// the migration SQL files.
type Config struct {
	// DatabaseURL is a Postgres connection URL, e.g.
	// "postgres://user:pass@host:5432/db?sslmode=disable"
	DatabaseURL string

	// MigrationsDir is the path to the directory containing *.sql migration
	// files (default: "migrations").
	MigrationsDir string
}

// Run applies all pending Up migrations using golang-migrate.
// The migrations table is named "schema_migrations".
func Run(cfg Config) error {
	dir := cfg.MigrationsDir
	if dir == "" {
		dir = "migrations"
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", dir),
		cfg.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("migrate: create instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate: up: %w", err)
	}
	return nil
}
