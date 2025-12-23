package database

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dsn string, log *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Error("database connect", err.Error())
		return nil, fmt.Errorf("database connect - %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Error("database connect", err.Error())
		return nil, fmt.Errorf("database connect - %v", err)
	}

	return db, nil
}

func Close(db *sqlx.DB, log *slog.Logger) {
	log.With(slog.String("operation", "close connection")).Info("closing the database connection")
	if err := db.Close(); err != nil {
		log.Error("closing database connection", err.Error())
	}
}
