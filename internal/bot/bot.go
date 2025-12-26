package bot

import (
	"context"
	"log/slog"

	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/bot/handlers"
	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Константы для описания операций
const (
	opStop      = "bot.StopListening"
	opListening = "bot.Listening"
)

// Updater описывает метод для обработки новых обновлений в канале
type Updater interface {
	HandleUpdate(ctx context.Context, update tgbotapi.Update) error
}

// Bot описывает API для взаимодействия с ботом
type Bot struct {
	log     *slog.Logger
	bot     *tgbotapi.BotAPI
	updater Updater
}

// NewBot конструктор для Bot
func NewBot(log *slog.Logger, token string, service *service.Service) (*Bot, error) {
	// Создаём новый экземпляр BotAPI по токену
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	// Инициализируем слой обработчиков с зависимостями
	h := handlers.NewHandler(log, service, bot)
	return &Bot{
		log:     log,
		bot:     bot,
		updater: h,
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
	b.log.Info("getting updates", slog.String("operation", opListening))
	// Читаем обновления в бесконечном цикле
	for update := range updates {
		// Вызываем обработчик для новых обновлений
		if err := b.updater.HandleUpdate(context.Background(), update); err != nil {
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

// Stop останавливает чтение обновлений из канала
func (b *Bot) Stop() {
	b.log.Info("stopping telegram bot", slog.String("operation", opStop))
	b.bot.StopReceivingUpdates()
}
