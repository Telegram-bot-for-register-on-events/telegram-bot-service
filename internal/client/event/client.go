package event

import (
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		log.Error("error connection to event gRPC-server", slog.String("address", address), slog.String("error", err.Error()))
		return nil, fmt.Errorf("error connection to event gRPC server: %w", err)
	}
	log.Info("connected to event gRPC server successfully")

	return &Client{
		log:    log,
		client: pb.NewEventServiceClient(conn),
		conn:   conn,
	}, nil
}
