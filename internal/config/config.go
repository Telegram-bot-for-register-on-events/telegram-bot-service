package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	DatabasePath     string
	GRPCPort         string
}

func getEnv(key, reserve string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return reserve
}

func LoadConfig(log *slog.Logger) (*Config, error) {
	log.Info("loading environment variables")
	if err := godotenv.Load(); err != nil {
		log.Error("load config", err.Error())
		return nil, fmt.Errorf("load config - %v", err)
	}

	return &Config{
		TelegramBotToken: getEnv("TG_BOT_TOKEN", ""),
		DatabasePath:     getEnv("DSN", ""),
		GRPCPort:         getEnv("GRPC_PORT", ""),
	}, nil
}
