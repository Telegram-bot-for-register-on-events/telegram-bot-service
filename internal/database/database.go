package database

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect открывает соединение с базой данной
func Connect(dsn string, log *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Error("error database connect", err.Error())
		return nil, fmt.Errorf("error database connect - %w", err)
	}

	// Проверяем подключение к базе данных, в противном случае возвращаем ошибку
	if err = db.Ping(); err != nil {
		log.Error("error database connect", err.Error())
		return nil, fmt.Errorf("error database connect - %w", err)
	}

	return db, nil
}

// Close закрывает соединение с базой данных
func Close(db *sqlx.DB, log *slog.Logger) {
	log.With(slog.String("operation", "close connection")).Info("closing the database connection")
	if err := db.Close(); err != nil {
		log.Error("closing database connection", err.Error())
	}
}
