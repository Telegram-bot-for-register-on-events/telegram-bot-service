package app

import (
	"log/slog"
	"os"

	grpcapp "github.com/Recrusion/telegram-bot-service/internal/app/grpc"
	bot "github.com/Recrusion/telegram-bot-service/internal/bot"
	"github.com/Recrusion/telegram-bot-service/internal/config"
	"github.com/Recrusion/telegram-bot-service/internal/database"
	"github.com/Recrusion/telegram-bot-service/internal/repository"
	"github.com/Recrusion/telegram-bot-service/internal/service"
	"github.com/jmoiron/sqlx"
)

// App описывает микросервис целиком, единая точка входа для всего микросервиса
type App struct {
	log        *slog.Logger
	GRPCServer *grpcapp.GRPCApp
	Bot        *bot.Bot
	Database   *sqlx.DB
}

// NewApp констурктор для App
func NewApp(log *slog.Logger, cfg *config.Config) *App {
	// Создаём gRPC-сервер приложения
	gRPCServer := grpcapp.NewApp(log, cfg.GRPCPort)
	log.Info("gRPC-server successfully created")

	// Устанавливаем подключение к базе данных
	db, err := database.Connect(cfg.DatabasePath, log)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	log.Info("connection to database successfully created")

	// Инициализируем слои (сервисный и для взаимодействия с базой данных)
	repo := repository.NewUserRepository(db, log)
	srvc := service.NewUserService(log, repo)

	// Создаём бота
	b, err := bot.NewBot(log, cfg.TelegramBotToken, srvc)
	if err != nil {
		log.Error("failed to create bot", "error", err)
		os.Exit(1)
	}
	log.Info("bot successfully created")

	return &App{
		log:        log,
		GRPCServer: gRPCServer,
		Bot:        b,
		Database:   db,
	}
}

// MustStart запускает приложение
// Бота и gRPC-сервер в отдельных горутинах
func (app *App) MustStart() {
	app.log.Info("application successfully started")
	go app.Bot.MustStart()
	go app.GRPCServer.MustStart()
}

// Close реализует GracefulShutdown для всего микросервиса
func (app *App) Stop() {
	app.log.Info("shutting down...")
	app.GRPCServer.Stop()
	app.Bot.Stop()
	database.Close(app.Database, app.log)
}
