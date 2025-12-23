package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Recrusion/telegram-bot-service/internal/grpc/event"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewApp(log *slog.Logger, port string) *App {
	gRPCServer := grpc.NewServer()
	event.Register(gRPCServer)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) start() error {
	log := a.log.With(
		slog.String("operation", "start gRPC server"),
		slog.String("port", a.port),
	)

	listener, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("start gRPC server - %v", err)
	}

	log.Info("starting gRPC server", slog.String("address", listener.Addr().String()))

	if err = a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("start gRPC server - %v", err)
	}

	return nil
}

func (a *App) MustStart() {
	if err := a.start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	a.log.With(slog.String("operation", "stop gRPC server")).Info("stopping gRPC server", slog.String("port", a.port))
	a.gRPCServer.GracefulStop()
}
