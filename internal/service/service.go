package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type UserService struct {
	log     *slog.Logger
	service UserSaver
}

type UserSaver interface {
	SaveUserInfo(ctx context.Context, chatID int64, username string) error
}

func NewUserService(log *slog.Logger, service UserSaver) *UserService {
	return &UserService{
		log:     log,
		service: service,
	}
}

func (s *UserService) SaveUserInfo(ctx context.Context, chatID int64, username string) error {
	if chatID == 0 {
		s.log.Error("")
		return errors.New("chatID cannot be equal to 0")
	} else if chatID < -999999999999999 || chatID > 999999999999999 {
		return errors.New("chatID out of range")
	}

	if username == "" {
		return errors.New("username cannot be empty")
	} else if len(username) < 5 || len(username) > 32 {
		return errors.New("username length must be between 5 and 32")
	}
	err := s.service.SaveUserInfo(ctx, chatID, username)
	if err != nil {
		return fmt.Errorf("save user info in service - %v", err)
	}
	return nil
}
