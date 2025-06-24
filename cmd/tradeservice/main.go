package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"tradeservice/internal/config"
	"tradeservice/internal/server/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DB struct {
	User   string
	Passwd string
	DBName string
	Host   string
	Port   string
}

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal("No config cannot start server", slog.Any("error", err))
	}

	sloger := logger.SetupLogger(cfg.Server.EnvType)
	sloger.Info("starting TradeService")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Passwd, cfg.DB.DBName)

	sloger.Info("Connecting db", psqlInfo, cfg)

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		fmt.Printf("Unable to acquire connection: %v\n", err)
		os.Exit(1)
	}

	println("Opening DB")

	sql, err := goose.OpenDBWithDriver("pgx", psqlInfo)
	if err != nil {
		log.Fatalf(err.Error())
	}

	println("Migrating")

	err = goose.Up(sql, "./internal/migrations")
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer db.Close()

}
