package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ChatID    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Save(ctx context.Context, chatID int64, username string) error {
	_, err := repo.db.NamedExecContext(ctx,
		"insert into users (chat_id, username, created_at) values (:chat_id, :username, :created_at) on conflict (chat_id) do nothing",
		User{
			ChatID:    chatID,
			Username:  username,
			CreatedAt: time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("ошибка сохранения данных о пользователе в базу данных: %v", err)
	}

	return nil
}
