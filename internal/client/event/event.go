package event

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
)

// Константы для описания операций
const (
	opGetEvents    = "event.GetEvents"
	opGetEvent     = "event.GetEvent"
	opRegisterUser = "event.RegisterUser"
)

// GetEvents метод для получения всех событий
func (c *Client) GetEvents(ctx context.Context) ([]*pb.Event, error) {
	// Отправляем запрос на другой микросервис
	response, err := c.client.GetEvents(ctx, &pb.GetEventsRequest{})
	if err != nil {
		c.log.Error("error", err.Error(), slog.String("operation", opGetEvents))
		return nil, fmt.Errorf("%s: %w", opGetEvents, err)
	}
	c.log.Info("getting events successfully", slog.Int("count", len(response.Events)), slog.String("operation", opGetEvents))
	return response.GetEvents(), nil
}

// GetEvent метод для получения конкретного события по его ID
func (c *Client) GetEvent(ctx context.Context, eventID string) (*pb.Event, error) {
	response, err := c.client.GetEvent(ctx, &pb.GetEventRequest{EventId: eventID})
	if err != nil {
		c.log.Error("error", err.Error(), slog.String("operation", opGetEvent))
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	c.log.Info("getting event successfully", slog.String("event_id", eventID), slog.String("operation", opGetEvent))
	return response.GetEvent(), nil
}

// RegisterUser метод для регистрации пользователя на конкретное событие
func (c *Client) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error) {
	response, err := c.client.RegisterUser(ctx, &pb.RegisterUserRequest{EventId: eventID, ChatId: chatID, Username: username})
	if err != nil {
		c.log.Error("error", err.Error(), slog.String("operation", opRegisterUser))
		return false, fmt.Errorf("%s: %w", opRegisterUser, err)
	}
	c.log.Info("register user on event successfully", slog.String("event_id", eventID), slog.String("username", username), slog.String("operation", opRegisterUser))
	return response.GetSuccess(), nil
}

// TODO: Добавить проверку и обработку на resp == nil и err == nil
