package storage

import (
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log/slog"
	"tradeservice/internal/storage/postgres"
)

func RunMigration(db *postgres.Storage, logger *slog.Logger, path string) error {

	logger.Info("Migrating")

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("couldn't setup migration %w", err)
	}

	sqlDb := stdlib.OpenDBFromPool(db.DB)
	if err := goose.Up(sqlDb, path); err != nil {
		return fmt.Errorf("couldn't run migration %w", err)
	}
	if err := sqlDb.Close(); err != nil {
		return fmt.Errorf("error creating migration %w", err)
	}

	return nil
}
