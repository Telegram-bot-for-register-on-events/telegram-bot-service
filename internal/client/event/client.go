package event

import (
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	opNewClient = "event.NewClient"
)

// Client описывает gRPC-клиент для взаимодействия с микросервисом событий
type Client struct {
	log    *slog.Logger
	client pb.EventServiceClient
	conn   *grpc.ClientConn
}

// NewClient конструктор для Client, устанавливает подключение к микросервису событий
func NewClient(log *slog.Logger, address string) (*Client, error) {
	// Устанавливаем gRPC-соединение
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opNewClient))
		return nil, fmt.Errorf("%s: %w", opNewClient, err)
	}

	return &Client{
		log:    log,
		client: pb.NewEventServiceClient(conn),
		conn:   conn,
	}, nil
}
