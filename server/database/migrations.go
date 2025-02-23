package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (d *Database) Migrate(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("vertica"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, d.Conn, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
