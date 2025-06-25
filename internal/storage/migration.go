package storage

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"log/slog"
)

func RunMigration(db *sql.DB, logger *slog.Logger) error {

	logger.Info("Migrating")

	err := goose.Up(db, "./internal/migrations")
	if err != nil {
		return err
	}

	return nil
}
