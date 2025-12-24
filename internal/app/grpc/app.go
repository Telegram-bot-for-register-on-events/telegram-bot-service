package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Recrusion/telegram-bot-service/internal/grpc/event"
	"google.golang.org/grpc"
)

// App описывает gRPC-сервер микросервиса
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

// NewApp конструктор для App
func NewApp(log *slog.Logger, port string) *App {
	gRPCServer := grpc.NewServer()
	event.Register(gRPCServer)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// start запускает gRPC-сервер на указанном порту для прослушивания входящих соединений и обработки запросов
func (a *App) start() error {
	log := a.log.With(
		slog.String("operation", "start gRPC server"),
		slog.String("port", a.port),
	)

	// Инициализируем слушателя на указанном порту, используя TCP-соединение
	listener, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		log.Error("start to gRPC server", err.Error())
		return fmt.Errorf("start gRPC server - %v", err)
	}

	log.Info("starting gRPC server", slog.String("address", listener.Addr().String()))

	// Принимаем входящие соединения от слушателя, создавая горутину для каждого из них
	// Cчитываем запросы и вызываем соответствующие обработчики для них
	if err = a.gRPCServer.Serve(listener); err != nil {
		log.Error("start to gRPC server", err.Error())
		return fmt.Errorf("start gRPC server - %v", err)
	}

	return nil
}

// MustStart обёртка для start, при ошибке - паникует
func (a *App) MustStart() {
	if err := a.start(); err != nil {
		panic(err)
	}
}

// Stop выполняет GracefulShutdown для gRPC-сервера
func (a *App) Stop() {
	a.log.With(slog.String("operation", "stop gRPC server")).Info("stopping gRPC server", slog.String("port", a.port))
	a.gRPCServer.GracefulStop()
}
