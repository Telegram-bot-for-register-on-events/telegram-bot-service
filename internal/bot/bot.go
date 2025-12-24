package adapter

import (
	"context"
	"log/slog"
	"time"

	"github.com/Recrusion/telegram-bot-service/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot описывает API для взаимодействия с ботом
type Bot struct {
	log     *slog.Logger
	bot     *tgbotapi.BotAPI
	service *service.UserService
}

// NewBot конструктор для Bot
func NewBot(log *slog.Logger, token string, service *service.UserService) (*Bot, error) {
	// Создаём новый экземпляр BotAPI по токену
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

// start начинает приём обновлений и прослушивание входящих сообщений из канала
func (b *Bot) start() error {
	// Получаем все обновления, начиная с самого первого
	u := tgbotapi.NewUpdate(0)
	// Устанавливаем тайм-аут, в течение которого будут прослушиваться входящие сообщения
	u.Timeout = 30
	// Инициализируем канал с обновлениями и устанавливаем долгоживущее подключение к серверам Telegram
	updates := b.bot.GetUpdatesChan(u)
	b.log.Info("getting updates", slog.String("operation", "listening incoming messages"))
	// Читаем обновления в бесконечном цикле
	for update := range updates {
		// Вызываем обработчик для новых обновлений
		if err := b.handleUpdate(context.Background(), update); err != nil {
			return err
		}
	}

	return nil
}

// MustStart обёртка для метода start, при ошибке - паникует
func (b *Bot) MustStart() {
	if err := b.start(); err != nil {
		panic(err)
	}
}

// handleUpdate принимает входящее обновление и вызывает обработчик для него
func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		return b.handleCommand(ctx, update.Message)
	}
	return nil
}

// handleCommand обработчик команд, принимает команду и вызывает для неё соответствующий обработчик
func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) error {
	b.log.Info("handling command", slog.String("command", msg.Command()), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Command() {
	// Вызов обработчика для команды "/start"
	case "start":
		b.log.Info("sending answer", slog.String("command", "start"), slog.Int("chat_id", int(msg.Chat.ID)))
		return b.startMessage(ctx, msg)
	}
	return nil
}

// startMessage обработчик для команды /start
func (b *Bot) startMessage(ctx context.Context, msg *tgbotapi.Message) error {
	// Формируем ответ для пользователя
	answer := tgbotapi.NewMessage(msg.Chat.ID,
		"Привет! 👋\nЯ бот для отслеживания и регистрации на события.",
	)

	b.log.Info("saving user info", slog.Int("chat_id", int(msg.Chat.ID)), slog.String("username", msg.From.UserName), slog.Time("created_at", time.Now()))
	// Сохраняем информацию в базу данных
	if err := b.service.SaveUserInfo(ctx, msg.Chat.ID, msg.From.UserName); err != nil {
		return err
	}

	// Отправляем ответ пользователю
	_, err := b.bot.Send(answer)
	if err != nil {
		b.log.Error("error send answer", err.Error())
		return err
	}
	return nil
}

// Stop останавливает чтение обновлений из канала
func (b *Bot) Stop() {
	b.log.With(slog.String("operation", "stop listening updates")).Info("stopping telegram bot")
	b.bot.StopReceivingUpdates()
}
