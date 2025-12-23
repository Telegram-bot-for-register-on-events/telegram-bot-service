package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/Recrusion/telegram-bot-service/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	repo *repository.UserRepository
}

func NewBot(token string, repo *repository.UserRepository) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Ошибка авторизации в бота: %v", err)
	}

	bot.Debug = true
	return &Bot{
		bot:  bot,
		repo: repo,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var msg tgbotapi.MessageConfig
			switch update.Message.Command() {
			case "start":
				err := b.repo.Save(context.Background(), update.Message.Chat.ID, update.Message.From.UserName)
				if err != nil {
					log.Printf("%v", err)
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! 👋\nЯ бот для отслеживания и регистрации на события.\nТвои данные сохранены.")
			}

			_, err := b.bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
