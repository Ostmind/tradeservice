package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"tradeservice/internal/config"
	categorieshandler "tradeservice/internal/server/handler/categories"
	productshandler "tradeservice/internal/server/handler/products"
	srv "tradeservice/internal/server/server"
	"tradeservice/internal/services/categories"
	"tradeservice/internal/services/product"
	"tradeservice/internal/storage"
	"tradeservice/internal/storage/postgres"
)

type App struct {
	server *srv.Server
	logger *slog.Logger
	db     *postgres.Storage
	cfg    *config.AppConfig
}

func New(logger *slog.Logger, cfg *config.AppConfig) (*App, error) {
	db, err := postgres.New(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish db connection %w", err)
	}

	categoryStorage, err := postgres.NewCategories(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't create categories %w", err)
	}

	productStorage, err := postgres.NewProducts(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't create products %w", err)
	}

	categoryManager := categories.New(categoryStorage)
	productManager := product.New(productStorage)

	categoryHandler := categorieshandler.NewCategoriesHandler(categoryManager, logger)
	productHandler := productshandler.NewProductHandler(productManager, logger)

	server := srv.New(logger, &cfg.Server, db, categoryHandler, productHandler)

	return &App{
		server: server,
		logger: logger,
		db:     db,
		cfg:    cfg,
	}, nil
}

func (a App) Run() {
	a.logger.Info("Starting app...")
	err := storage.RunMigration(a.db, a.logger, a.cfg.Server.MigrationPath)

	if err != nil {
		a.logger.Error("couldn't run migrations %w", slog.Any("error_details", err))
	}

	a.server.Run()
}

func (a App) Stop(ctx context.Context, shutdownTimeout time.Duration) {
	a.logger.Info("Stopping app...")

	timeout := shutdownTimeout

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	doneCh := make(chan error)
	go func() {
		doneCh <- a.server.Stop(ctxWithTimeout)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			a.logger.Error("Error while stopping server: %v", slog.Any("error_details", err))
		}

		a.logger.Info("App has been stopped gracefully")

	case <-ctx.Done():
		a.logger.Warn("App stopped forced")
	}
}
