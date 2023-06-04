package database

import (
	"database/sql"
	"fmt"

	"github.com/andredubov/todo-backend/internal/config"
)

func NewPostgresConnection(cfg config.Config) (*sql.DB, error) {

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.DatabaseName, cfg.Postgres.SSLMode, cfg.Postgres.Password))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
