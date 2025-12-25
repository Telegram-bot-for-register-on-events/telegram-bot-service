package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config описывает конфигурацию микросервиса
type Config struct {
	telegramBotConfig *telegramBotConfig
	databaseConfig    *databaseConfig
	gRPCClientConfig  *gRPCClientConfig
}

// telegramBotConfig описывает конфигурацию телеграм-бота
type telegramBotConfig struct {
	token string
}

// databaseConfig описывает конфигурацию базы данных
type databaseConfig struct {
	path string
}

// gRPCClientConfig описывает конфигурацию gRPC-клиента
type gRPCClientConfig struct {
	address string
}

// newTelegramBotConfig загружает конфигурацию для телеграм-бота
func newTelegramBotConfig(log *slog.Logger) (*telegramBotConfig, error) {
	token := getEnv("TELEGRAM_BOT_TOKEN", "")
	if token == "" {
		log.Error("telegram-bot token cannot be empty")
		return nil, errors.New("telegram-bot token cannot be empty")
	}
	tgBotCfg := &telegramBotConfig{token: token}
	return tgBotCfg, nil
}

// newDatabaseConfig загружает конфигурацию для базы данных
func newDatabaseConfig(log *slog.Logger) (*databaseConfig, error) {
	path := getEnv("DSN", "")
	if path == "" {
		log.Error("dsn cannot be empty")
		return nil, errors.New("dsn cannot be empty")
	}
	dbCfg := &databaseConfig{path: path}
	return dbCfg, nil
}

// newGRPCClientConfig загружает конфигурацию для gRPC-клиента
func newGRPCClientConfig(log *slog.Logger) (*gRPCClientConfig, error) {
	address := getEnv("GRPC_ADDRESS", "")
	if address == "" {
		log.Error("gRPC address cannot be empty")
		return nil, errors.New("gRPC address cannot be empty")
	}
	gRPCCfg := &gRPCClientConfig{address: address}
	return gRPCCfg, nil
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
		log.Error("error load config", err.Error())
		return nil, fmt.Errorf("error load config - %w", err)
	}
	log.Info("environment variables successfully loaded")

	// Загружаем конфигурацию базы данных
	dbCfg, err := newDatabaseConfig(log)
	if err != nil {
		log.Error("error load config", err.Error())
		return nil, err
	}
	// Загружаем конфигурацию телеграм-бота
	tgBotCfg, err := newTelegramBotConfig(log)
	if err != nil {
		log.Error("error load config", err.Error())
		return nil, err
	}
	// Загружаем конфигурацию gRPC-клиента
	gRPCCfg, err := newGRPCClientConfig(log)
	if err != nil {
		log.Error("error load config", err.Error())
		return nil, err
	}

	return &Config{
		telegramBotConfig: tgBotCfg,
		databaseConfig:    dbCfg,
		gRPCClientConfig:  gRPCCfg,
	}, nil
}

// MustLoadConfig обёртка для LoadConfig, при ошибке - паникует
func MustLoadConfig(log *slog.Logger) *Config {
	cfg, err := LoadConfig(log)
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetTelegramBotToken геттер, для получения значения токена телеграм-бота
func (c *Config) GetTelegramBotToken() string {
	return c.telegramBotConfig.token
}

// GetDatabasePath геттер, для получения пути подключения к базе данных
func (c *Config) GetDatabasePath() string {
	return c.databaseConfig.path
}

// GetGRPCAddress геттер, для получения адреса gRPC-сервера
func (c *Config) GetGRPCAddress() string {
	return c.gRPCClientConfig.address
}
