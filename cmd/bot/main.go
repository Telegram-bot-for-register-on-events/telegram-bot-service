package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Recrusion/telegram-bot-service/internal/app"
	"github.com/Recrusion/telegram-bot-service/internal/config"
)

func main() {
	// Инициализируем логгер
	log := setupLogger()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	log.Info("config successfully loaded")

	// Создаём новый инстанс микросервиса
	application := app.NewApp(log, cfg)
	// Запускаем его
	application.MustStart()

	// Создаём канал для приёма сигналов операционной системы
	stop := make(chan os.Signal, 1)
	// Передаём входящие сигналы в канал stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	// Читаем из канала, пока не придёт соответствующий сигнал
	<-stop

}

// setupLogger инициализирует логгер с JSON-обработчиком
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
