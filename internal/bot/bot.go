package adapter

import (
	"context"
	"log/slog"
	"time"

	"github.com/Recrusion/telegram-bot-service/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	log     *slog.Logger
	bot     *tgbotapi.BotAPI
	service *service.UserService
}

func NewBot(log *slog.Logger, token string, service *service.UserService) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		log:     log,
		bot:     bot,
		service: service,
	}, nil
}

func (b *Bot) start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.bot.GetUpdatesChan(u)
	b.log.Info("getting updates", slog.String("operation", "listening incoming messages"))
	for update := range updates {
		if err := b.handleUpdate(context.Background(), update); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) MustStart() {
	if err := b.start(); err != nil {
		panic(err)
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		return b.handleCommand(ctx, update.Message)
	}
	return nil
}

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) error {
	b.log.Info("handling command", slog.String("command", msg.Command()), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Command() {
	case "start":
		b.log.Info("sending answer", slog.String("command", "start"), slog.Int("chat_id", int(msg.Chat.ID)))
		return b.startMessage(ctx, msg)
	}
	return nil
}

func (b *Bot) startMessage(ctx context.Context, msg *tgbotapi.Message) error {
	answer := tgbotapi.NewMessage(msg.Chat.ID,
		"Привет! 👋\nЯ бот для отслеживания и регистрации на события.",
	)

	b.log.Info("saving user info", slog.Int("chat_id", int(msg.Chat.ID)), slog.String("username", msg.From.UserName), slog.Time("created_at", time.Now()))
	if err := b.service.SaveUserInfo(ctx, msg.Chat.ID, msg.From.UserName); err != nil {
		return err
	}

	_, err := b.bot.Send(answer)
	if err != nil {
		return nil
	}
	return nil
}

func (b *Bot) Stop() {
	b.log.With(slog.String("operation", "stop listening updates")).Info("stopping telegram bot")
	b.bot.StopReceivingUpdates()
}
