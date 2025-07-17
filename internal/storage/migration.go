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
		panic(err)
	}

	sqlDb := stdlib.OpenDBFromPool(db.DB)
	if err := goose.Up(sqlDb, path); err != nil {
		panic(err)
	}
	if err := sqlDb.Close(); err != nil {
		return fmt.Errorf("error creating migration %w", err)
	}

	return nil
}
