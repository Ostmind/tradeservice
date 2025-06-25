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
	defer db.Close()

	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}
