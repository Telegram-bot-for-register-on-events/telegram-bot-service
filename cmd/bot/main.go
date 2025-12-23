package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Recrusion/telegram-bot-service/internal/app"
	bot "github.com/Recrusion/telegram-bot-service/internal/bot"
	"github.com/Recrusion/telegram-bot-service/internal/config"
	"github.com/Recrusion/telegram-bot-service/internal/database"
	"github.com/Recrusion/telegram-bot-service/internal/repository"
	"github.com/Recrusion/telegram-bot-service/internal/service"
)

func main() {
	log := setupLogger()

	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	log.Info("config successfully loaded")

	db, err := database.Connect(cfg.DatabasePath, log)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	log.Info("connection to database successfully created")
	defer database.Close(db, log)

	repo := repository.NewUserRepository(db, log)
	srvc := service.NewUserService(log, repo)

	b, err := bot.NewBot(log, cfg.TelegramBotToken, srvc)
	if err != nil {
		panic(err)
	}
	application := app.NewApp(log, cfg.GRPCPort)

	go b.MustStart()
	go application.GRPCServer.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down...")
	application.GRPCServer.Stop()

	b.Stop()
}

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
