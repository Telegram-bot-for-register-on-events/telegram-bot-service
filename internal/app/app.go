package app

import (
	"log/slog"
	"os"

	"github.com/Recrusion/telegram-bot-service/internal/bot"
	"github.com/Recrusion/telegram-bot-service/internal/client/event"
	"github.com/Recrusion/telegram-bot-service/internal/config"
	"github.com/Recrusion/telegram-bot-service/internal/database"
	"github.com/Recrusion/telegram-bot-service/internal/repository"
	"github.com/Recrusion/telegram-bot-service/internal/service"
	"github.com/jmoiron/sqlx"
)

// App описывает микросервис целиком, единая точка входа для всего микросервиса
type App struct {
	log      *slog.Logger
	Bot      *bot.Bot
	Database *sqlx.DB
	Client   *event.Client
}

// NewApp конструктор для App
func NewApp(log *slog.Logger) *App {
	// Инициализируем конфиг
	cfg := newCfg(log)
	// Создаём gRPC-клиент для отправки запросов
	client := newClient(log, cfg)
	// Создаём подключение к базе данных
	db := dbConn(log, cfg)
	// Инициализируем слои (сервисный и для взаимодействия с базой данных)
	repo := repository.NewUserRepository(db, log)
	srvc := service.NewUserService(log, client, client, repo)

	b := newBot(log, cfg, srvc)
	return &App{
		log:      log,
		Bot:      b,
		Database: db,
		Client:   client,
	}
}

// MustStart запускает приложение
func (app *App) MustStart() {
	app.log.Info("application successfully started")
	go app.Bot.MustStart()
}

// Stop реализует GracefulShutdown для всего микросервиса
func (app *App) Stop() {
	app.log.Info("shutting down...")
	app.Bot.Stop()
	database.Close(app.Database, app.log)
}

// newCfg обёртка для инициализации объекта конфигурации
func newCfg(log *slog.Logger) *config.Config {
	cfg := config.MustLoadConfig(log)
	log.Info("config successfully loaded")
	return cfg
}

// dbConn обёртка для установки соединения к базе данных
func dbConn(log *slog.Logger, cfg *config.Config) *sqlx.DB {
	db, err := database.Connect(cfg.GetDatabasePath(), log)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	log.Info("connection to database successfully created")
	return db
}

// newBot обёртка для создания нового экземпляра BotAPI по токену
func newBot(log *slog.Logger, cfg *config.Config, srvc *service.UserService) *bot.Bot {
	b, err := bot.NewBot(log, cfg.GetTelegramBotToken(), srvc)
	if err != nil {
		log.Error("failed to create bot", "error", err)
		os.Exit(1)
	}
	log.Info("bot successfully created")
	return b
}

// newClient обёртка для создания gRPC-клиента
func newClient(log *slog.Logger, cfg *config.Config) *event.Client {
	client, err := event.NewClient(log, cfg.GetGRPCAddress())
	if err != nil {
		log.Error("failed to create client", "error", err)
		os.Exit(1)
	}
	log.Info("client successfully created")
	return client
}
