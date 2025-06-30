package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"tradeservice/internal/config"
)

type Storage struct {
	DB *pgxpool.Pool
}

func New(dbConfig config.DBConfig) (*Storage, error) {
	db := &Storage{}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Passwd, dbConfig.DBName)

	err := db.connect(psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("error creating connection DB %w", err)
	}
	return db, nil
}

func (store *Storage) Close() {
	store.DB.Close()
}

func (store *Storage) connect(connStr string) error {
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return fmt.Errorf("db.connect: %v", err)
	}

	err = pool.Ping(context.Background())

	if err != nil {
		return fmt.Errorf("db.connect pool ping: %v", err)
	}

	store.DB = pool
	return nil
}
