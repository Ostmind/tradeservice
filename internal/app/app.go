package app

import (
	"context"
	"log"
	"log/slog"
	"time"
	"tradeservice/internal/config"
	"tradeservice/internal/server/handler"
	srv "tradeservice/internal/server/server"
	"tradeservice/internal/services/categories"
	"tradeservice/internal/services/product"
	"tradeservice/internal/storage/postgres"
)

type App struct {
	server *srv.Server
	logger *slog.Logger
	db     *postgres.Storage
}

func New(logger *slog.Logger, cfg *config.AppConfig) *App {

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("couldn't establish db connection %w", err)
	}

	categoryStorage, err := postgres.NewCategories(db)
	if err != nil {
		log.Fatalf("couldn't create categories %w", err)
	}

	productStorage, err := postgres.NewProducts(db)
	if err != nil {
		log.Fatalf("couldn't create products %w", err)
	}

	categoryManager := categories.New(categoryStorage)
	productManager := product.New(productStorage)

	categoryHandler := handlers.NewCategoriesHandler(categoryManager, logger)
	productHandler := handlers.NewProductHandler(productManager, logger)

	server := srv.New(logger, &cfg.Server, db, categoryHandler, productHandler)

	return &App{
		server: server,
		logger: logger,
		db:     db,
	}
}

func (a App) Run() {
	a.logger.Info("Starting app...")
	go a.server.Run()
}

func (a App) Stop(ctx context.Context, shutdownTimeout time.Duration) {
	a.logger.Info("Stopping app...")

	timeout := shutdownTimeout * time.Second
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	doneCh := make(chan error)
	go func() {
		doneCh <- a.server.Stop(ctxWithTimeout)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			a.logger.Error("Error while stopping server: %v", err)
		}
		a.logger.Info("App has been stopped gracefully")

	case <-ctx.Done():
		a.logger.Warn("App stopped forced")
	}
}
