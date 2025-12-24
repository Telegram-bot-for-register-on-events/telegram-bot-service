package event

import (
	"context"

	"github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
)

// Event определяет методы объявленные в event.proto
type Event interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, eventID string) (Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
}

// serverAPI описывает gRPC-сервер
type serverAPI struct {
	event.UnimplementedEventServiceServer // Базовая реализация из сгенерированного кода
	event                                 Event
}

// Register добавляет event сервис в общий gRPC-сервер
func Register(gRPC *grpc.Server) {
	event.RegisterEventServiceServer(gRPC, &serverAPI{})
}

// GetEvents метод для получения всех событий
func (s *serverAPI) GetEvents(ctx context.Context, req *event.GetEventsRequest) (*event.GetEventsResponse, error) {
	panic("implement me")
}

// GetEvent метод для получения конкретного события по его ID
func (s *serverAPI) GetEvent(ctx context.Context, req *event.GetEventRequest) (*event.GetEventResponse, error) {
	panic("implement me")
}

// RegisterUser метод для регистрации пользователя на конкретное событие
func (s *serverAPI) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	panic("implement me")
}
