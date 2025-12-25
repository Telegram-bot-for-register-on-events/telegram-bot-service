package bot

import (
	"context"
	"log/slog"
	"time"

	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/client/event"
	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	opStop = "stop listening updates"
)

// Bot описывает API для взаимодействия с ботом
type Bot struct {
	log     *slog.Logger
	bot     *tgbotapi.BotAPI
	service *service.UserService
	client  *event.Client
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
		return b.handleMessage(ctx, update.Message)
	}
	return nil
}

// handleMessage обработчик сообщений, проверяет, что пришло: команда или текст, и на основе этого вызывает соответствующий обработчик
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	b.log.Info("handling message", slog.String("text", msg.Text), slog.Int("chat_id", int(msg.Chat.ID)))
	if msg.IsCommand() {
		return b.handleCommand(ctx, msg)
	}
	if msg.Text != "" {
		return b.handleText(ctx, msg)
	}
	return nil
}

// handleCommand метод для обработки команд
func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) error {
	b.log.Info("handling command", slog.String("command", msg.Command()), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Command() {
	case "start":
		return b.startMessage(ctx, msg)
	case "getEvents":
		return b.showEvents(ctx, msg)
	}
	return nil
}

// handleText метод для обработки сообщений
func (b *Bot) handleText(ctx context.Context, msg *tgbotapi.Message) error {
	b.log.Info("handling text", slog.String("text", msg.Text), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Text {
	case "Посмотреть предстоящие события":
		return b.showEvents(ctx, msg)
	}
	return nil
}

// startMessage обработчик для команды /start
func (b *Bot) startMessage(ctx context.Context, msg *tgbotapi.Message) error {
	// Формируем ответ для пользователя
	answer := tgbotapi.NewMessage(msg.Chat.ID,
		"Привет! 👋\nЯ бот для отслеживания и регистрации на события.",
	)
	// После приветствия показываем основную клавиатуру
	answer.ReplyMarkup = mainKeyboard()

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

// showEvents обработчик для команды /getEvents
func (b *Bot) showEvents(ctx context.Context, msg *tgbotapi.Message) error {
	// Отправляем данные в сервисный слой, в случае ошибки - отправляем пользователю соответствующее сообщение
	events, err := b.service.GetEvents(ctx)
	if err != nil {
		errMsg := tgbotapi.NewMessage(msg.Chat.ID, "Произошла ошибка")
		b.bot.Send(errMsg)
		return err
	}

	if len(events) == 0 {
		noEventsMsg := tgbotapi.NewMessage(msg.Chat.ID, "К сожалению, событий не найдено")
		noEventsMsg.ReplyMarkup = mainKeyboard()
		b.bot.Send(noEventsMsg)
		return nil
	}

	// Создаём "кнопки" с соответствующими данными
	var eventsButtons []EventButton
	for _, event := range events {
		eventsButtons = append(eventsButtons, EventButton{
			EventID: event.Id,
			Title:   event.Title,
		})
	}

	message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите событие, для просмотра детальной информации")
	message.ParseMode = tgbotapi.ModeMarkdown
	// Отправляем клавиатуру с событиями пользователю
	message.ReplyMarkup = eventsKeyboard(eventsButtons)
	_, err = b.bot.Send(message)
	if err != nil {
		b.log.Error("error send answer", err.Error())
		return err
	}
	return nil
}

// Stop останавливает чтение обновлений из канала
func (b *Bot) Stop() {
	b.log.With(slog.String("operation", opStop)).Info("stopping telegram bot")
	b.bot.StopReceivingUpdates()
}
