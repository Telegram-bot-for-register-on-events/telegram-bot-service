package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
)

// Константы для описания операций
const (
	opSaveUserInfo = "service.SaveUserInfo"
	opGetEvents    = "service.GetEvents"
	opGetEvent     = "service.GetEvent"
	opRegisterUser = "service.RegisterUser"
)

// Service описывает сервисный слой микросервиса
type Service struct {
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

// NewService конструктор для создания Service
func NewService(log *slog.Logger, eventReceiver EventReceiver, userRegister UserRegister, userSaver UserSaver) *Service {
	return &Service{
		log:           log,
		eventReceiver: eventReceiver,
		userRegister:  userRegister,
		userSaver:     userSaver,
	}
}

// SaveUserInfo проводит валидацию входных данных и передаёт их в слой взаимодействия с базой данных
func (s *Service) SaveUserInfo(ctx context.Context, chatID int64, username string) error {
	if err := validateChatID(chatID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return err
	}

	if err := validateUsername(username); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return err
	}

	err := s.userSaver.SaveUserInfo(ctx, chatID, username)
	if err != nil {
		return fmt.Errorf("%s: %w", opSaveUserInfo, err)
	}
	return nil
}

// GetEvents отправляет данные для получения всех событий
func (s *Service) GetEvents(ctx context.Context) ([]*pb.Event, error) {
	events, err := s.eventReceiver.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetEvents, err)
	}
	return events, nil
}

// GetEvent проводит валидацию входных данных и отправляет их для получения конкретного события
func (s *Service) GetEvent(ctx context.Context, eventID string) (*pb.Event, error) {
	if err := validateEventID(eventID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opGetEvent))
		return nil, err
	}

	event, err := s.eventReceiver.GetEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	return event, nil
}

// RegisterUser валидирует входные данные и отправляет их для регистрации пользователя на конкретное событие
func (s *Service) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error) {
	if err := validateEventID(eventID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return false, err
	}

	if err := validateUsername(username); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return false, err
	}

	if err := validateChatID(chatID); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return false, err
	}

	result, err := s.userRegister.RegisterUser(ctx, eventID, chatID, username)
	if err != nil {
		return false, fmt.Errorf("%s: %w", opRegisterUser, err)
	}
	return result, nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	} else if len(username) < 5 || len(username) > 32 {
		return errors.New("username length must be between 5 and 32")
	}
	return nil
}

func validateChatID(chatID int64) error {
	if chatID == 0 {
		return errors.New("chatID cannot be equal to 0")
	} else if chatID < -999999999999999 || chatID > 999999999999999 {
		return errors.New("chatID out of range")
	}
	return nil
}

func validateEventID(eventID string) error {
	if eventID == "" {
		return errors.New("eventID cannot be empty")
	}
	return nil
}
