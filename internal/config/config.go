package config

import (
	"fmt"
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

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load config - %v", err)
	}

	return &Config{
		TelegramBotToken: getEnv("TG_BOT_TOKEN", ""),
		DatabasePath:     getEnv("DSN", ""),
		GRPCPort:         getEnv("GRPC_PORT", ""),
	}, nil
}
