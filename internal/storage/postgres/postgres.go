package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"tradeservice/internal/config"
)

type Storage struct {
	DB *sql.DB
}

func New(dbConfig config.DBConfig) (*Storage, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Passwd, dbConfig.DBName)

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("error Creating Connection DB %s", err)
	}
	return &Storage{DB: db}, nil
}

func (store Storage) Close() error {
	err := store.DB.Close()
	return fmt.Errorf("error Shuttingdown DB %s", err)
}
