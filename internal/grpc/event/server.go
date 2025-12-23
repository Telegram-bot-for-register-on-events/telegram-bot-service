package event

import (
	"context"

	event "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
)

type Event interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, eventID string) (Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
}

type serverAPI struct {
	event.UnimplementedEventServiceServer
	event Event
}

func Register(gRPC *grpc.Server) {
	event.RegisterEventServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) GetEvents(ctx context.Context, req *event.GetEventsRequest) (*event.GetEventsResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetEvent(ctx context.Context, req *event.GetEventRequest) (*event.GetEventResponse, error) {
	panic("implement me")
}

func (s *serverAPI) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	panic("implement me")
}
