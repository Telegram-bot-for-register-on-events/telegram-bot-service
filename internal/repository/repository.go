package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

// User описывает данные о пользовователе, необходимые для сохранения
type User struct {
	ChatID    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

// UserRepository описывает объект базы данных
type UserRepository struct {
	log *slog.Logger
	db  *sqlx.DB
}

// NewUserRepository конструктор для создания UserRepository
func NewUserRepository(db *sqlx.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

// SaveUserInfo метод для сохранения информации в базе данных
func (repo *UserRepository) SaveUserInfo(ctx context.Context, chatID int64, username string) error {
	// Выполняем INSERT-запрос
	_, err := repo.db.NamedExecContext(ctx,
		"insert into users (chat_id, username, created_at) values (:chat_id, :username, :created_at) on conflict (chat_id) do nothing",
		User{
			ChatID:    chatID,
			Username:  username,
			CreatedAt: time.Now(),
		},
	)

	if err != nil {
		repo.log.Error("save user info", err.Error())
		return fmt.Errorf("save user info in repo - %v", err)
	}

	return nil
}
