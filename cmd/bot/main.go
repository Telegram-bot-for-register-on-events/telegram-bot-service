package main

import (
	"log"

	bot2 "github.com/Recrusion/telegram-bot-service/internal/bot"
	"github.com/Recrusion/telegram-bot-service/internal/config"
	"github.com/Recrusion/telegram-bot-service/internal/database"
	"github.com/Recrusion/telegram-bot-service/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Конфиг успешно загружен!")

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Подключение к базе данных установлено!")
	defer db.Close()

	repo := repository.NewUserRepository(db)
	bot, err := bot2.NewBot(cfg.TelegramBotToken, repo)
	if err != nil {
		log.Panic(err)
	}

	bot.Start()
}
