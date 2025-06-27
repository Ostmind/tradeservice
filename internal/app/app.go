package app

import (
	"context"
	"log"
	"log/slog"
	"time"
	"tradeservice/internal/config"
	srv "tradeservice/internal/server/server"
	"tradeservice/internal/storage/postgres"
)

type App struct {
	server *srv.Server
	logger *slog.Logger
	db     *postgres.Storage
}

func New(logger *slog.Logger, cfg *config.AppConfig) *App {

	//handler := handl.New(logger, client)

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("Couldn't establish db connection %s", err)
	}

	server := srv.New(logger, &cfg.Server, db)

	log.Print("Config: ", cfg)

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
