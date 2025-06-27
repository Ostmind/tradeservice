package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tradeservice/internal/app"
	"tradeservice/internal/config"
	"tradeservice/internal/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("No config cannot start server", slog.Any("error", err))
	}

	sloger := logger.SetupLogger(cfg.Server.EnvType)
	sloger.Info("starting TradeService")

	app := app.New(sloger, cfg)
	app.Run()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)

	<-stopChan
	sloger.Info("Recieved interrupt signal")
	app.Stop(context.Background(), cfg.Server.ShutdownTimeout)
}
