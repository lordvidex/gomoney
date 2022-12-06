package postgres

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/config"
	"github.com/pkg/errors"
)

// NewConn creates a new database connection.
func NewConn(ctx context.Context, c *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, c.Get("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}
	err = runMigrations(c)
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return conn, nil
}

func runMigrations(c *config.Config) error {
	dr := c.Get("MIGRATION_DIRECTORY")
	if dr == "" {
		dr = "file://server/internal/adapters/postgres/migrations"
	}
	m, err := migrate.New(dr, c.Get("DATABASE_URL"))
	if err != nil {
		return errors.Wrap(err, "failed to run migrations")
	}
	return m.Up()
}
