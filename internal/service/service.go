package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
)

// UserService описывает сервисный слой микросервиса
type UserService struct {
	log           *slog.Logger
	eventReceiver EventReceiver
	userRegister  UserRegister
	userSaver     UserSaver
}

// EventReceiver описывает методы для получения информации о событиях
type EventReceiver interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
}

// UserRegister описывает метод для регистрации пользователя на конкретное событие
type UserRegister interface {
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
}

// UserSaver определяет методы для сохранения информации о пользователе
type UserSaver interface {
	SaveUserInfo(ctx context.Context, chatID int64, username string) error
}

// NewUserService конструктор для создания UserService
func NewUserService(log *slog.Logger, eventReceiver EventReceiver, userRegister UserRegister, userSaver UserSaver) *UserService {
	return &UserService{
		log:           log,
		eventReceiver: eventReceiver,
		userRegister:  userRegister,
		userSaver:     userSaver,
	}
}

// SaveUserInfo проводит валидацию входных данных и передаёт их в слой взаимодействия с базой данных
func (s *UserService) SaveUserInfo(ctx context.Context, chatID int64, username string) error {
	if chatID == 0 {
		s.log.Error("chatID cannot be equal to 0")
		return errors.New("chatID cannot be equal to 0")
	} else if chatID < -999999999999999 || chatID > 999999999999999 {
		return errors.New("chatID out of range")
	}

	if username == "" {
		return errors.New("username cannot be empty")
	} else if len(username) < 5 || len(username) > 32 {
		return errors.New("username length must be between 5 and 32")
	}

	err := s.userSaver.SaveUserInfo(ctx, chatID, username)
	if err != nil {
		return fmt.Errorf("error save user info in service - %w", err)
	}
	return nil
}

// GetEvents отправляет данные для получения всех событий
func (s *UserService) GetEvents(ctx context.Context) ([]*pb.Event, error) {
	events, err := s.eventReceiver.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting events - %w", err)
	}
	return events, nil
}

// GetEvent проводит валидацию входных данных и отправляет их для получения конкретного события
func (s *UserService) GetEvent(ctx context.Context, eventID string) (*pb.Event, error) {
	if eventID == "" {
		s.log.Error("eventID cannot be empty")
		return nil, errors.New("eventID cannot be empty")
	}

	event, err := s.eventReceiver.GetEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("error getting event - %w", err)
	}
	return event, nil
}

// RegisterUser валидирует входные данные и отправляет их для регистрации пользователя на конкретное событие
func (s *UserService) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error) {
	if eventID == "" {
		s.log.Error("eventID cannot be empty")
		return false, errors.New("eventID cannot be empty")
	}

	if chatID == 0 {
		s.log.Error("chatID cannot be equal to 0")
		return false, errors.New("chatID cannot be equal to 0")
	} else if chatID < -999999999999999 || chatID > 999999999999999 {
		return false, errors.New("chatID out of range")
	}

	if username == "" {
		return false, errors.New("username cannot be empty")
	} else if len(username) < 5 || len(username) > 32 {
		return false, errors.New("username length must be between 5 and 32")
	}

	result, err := s.userRegister.RegisterUser(ctx, eventID, chatID, username)
	if err != nil {
		return false, fmt.Errorf("error registering user - %w", err)
	}
	return result, nil
}
