package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config описывает конфигурацию микросервиса
type Config struct {
	TelegramBotToken string
	DatabasePath     string
	GRPCPort         string
}

// getEnv проверяет наличие переменной окружения и возвращает её текущее значение, либо стандартное, при отсутствии текущего
func getEnv(key, reserve string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return reserve
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig(log *slog.Logger) (*Config, error) {
	log.Info("loading environment variables")
	// Чтение переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Error("load config", err.Error())
		return nil, fmt.Errorf("load config - %v", err)
	}
	log.Info("environment variables successfully loaded")
	
	return &Config{
		TelegramBotToken: getEnv("TG_BOT_TOKEN", ""),
		DatabasePath:     getEnv("DSN", ""),
		GRPCPort:         getEnv("GRPC_PORT", ""),
	}, nil
}
