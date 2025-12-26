package database

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Константы для описания операций
const (
	opConnect         = "database.connect"
	opCloseConnection = "database.closeConnection"
)

// Connect открывает соединение с базой данной
func Connect(driverName, dsn string, log *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Error("operation", opConnect, err.Error())
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	// Проверяем подключение к базе данных, в противном случае возвращаем ошибку
	if err = db.Ping(); err != nil {
		log.Error("operation", opConnect, err.Error())
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	return db, nil
}

// Close закрывает соединение с базой данных
func Close(db *sqlx.DB, log *slog.Logger) {
	log.With(slog.String("operation", opCloseConnection)).Info("closing the database connection")
	if err := db.Close(); err != nil {
		log.Error("closing database connection", err.Error())
	}
}
