package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Константы для описания операций
const (
	opSaveUserInfo = "repo.SaveUserInfo"
)

// User описывает данные о пользователе, необходимые для сохранения
type User struct {
	ChatID    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

// Storage описывает объект базы данных
type Storage struct {
	log *slog.Logger
	DB  *sqlx.DB
}

// Константы для описания операций
const (
	opConnect         = "postgres.connect"
	opCloseConnection = "postgres.closeConnection"
)

// NewStorage устанавливает соединение с базой данных, конструктор для Storage
func NewStorage(log *slog.Logger, driverName, dsn string) (*Storage, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConnect))
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	// Проверяем подключение к базе данных, в противном случае возвращаем ошибку
	if err = db.Ping(); err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConnect))
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	return &Storage{
		DB:  db,
		log: log,
	}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() {
	s.log.Info("close db connection..", slog.String("operation", opCloseConnection))
	if err := s.DB.Close(); err != nil {
		s.log.Error("closing database connection", slog.String("error", err.Error()))
	}
}

// SaveUserInfo метод для сохранения информации в базе данных
func (s *Storage) SaveUserInfo(ctx context.Context, chatID int64, username string) error {
	// Выполняем INSERT-запрос
	_, err := s.DB.NamedExecContext(ctx,
		"insert into users (chat_id, username, created_at) values (:chat_id, :username, :created_at) on conflict (chat_id) do nothing",
		User{
			ChatID:    chatID,
			Username:  username,
			CreatedAt: time.Now(),
		},
	)

	if err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opSaveUserInfo))
		return fmt.Errorf("%s: %w", opSaveUserInfo, err)
	}

	return nil
}
