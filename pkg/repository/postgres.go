package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		log.Fatalf("error open db: %s", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error connection to DB: %s", err)
		return nil, err
	}

	return db, nil
}
