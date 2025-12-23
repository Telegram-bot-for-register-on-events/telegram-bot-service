package app

import (
	"log/slog"

	grpcapp "github.com/Recrusion/telegram-bot-service/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, gRPCPort string) *App {
	gRPCServer := grpcapp.NewApp(log, gRPCPort)
	return &App{GRPCServer: gRPCServer}
}
