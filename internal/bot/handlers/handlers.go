package handlers

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/bot/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service описывает методы для взаимодействия с сервисным слоем
type Service interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
	SaveUserInfo(ctx context.Context, chatID int64, username string) error
}

// Sender описывает метод для отправки сообщений пользователю
type Sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// Handler описывает слой обработчиков для телеграм-бота
type Handler struct {
	log     *slog.Logger
	service Service
	sender  Sender
}

// NewHandler конструктор для Handler
func NewHandler(log *slog.Logger, service Service, sender Sender) *Handler {
	return &Handler{
		log:     log,
		service: service,
		sender:  sender,
	}
}

// HandleUpdate принимает входящее обновление и вызывает обработчик для него
func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}
	return nil
}

// handleMessage обработчик сообщений, проверяет, что пришло: команда или текст, и на основе этого вызывает соответствующий обработчик
func (h *Handler) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	h.log.Info("handling message", slog.String("text", msg.Text), slog.Int("chat_id", int(msg.Chat.ID)))
	if msg.IsCommand() {
		return h.handleCommand(ctx, msg)
	}
	if msg.Text != "" {
		return h.handleText(ctx, msg)
	}
	return nil
}

// handleCommand метод для обработки команд
func (h *Handler) handleCommand(ctx context.Context, msg *tgbotapi.Message) error {
	h.log.Info("handling command", slog.String("command", msg.Command()), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Command() {
	case "start":
		return h.startMessage(ctx, msg)
	case "getEvents":
		return h.showEvents(ctx, msg)
	}
	return nil
}

// handleText метод для обработки сообщений
func (h *Handler) handleText(ctx context.Context, msg *tgbotapi.Message) error {
	h.log.Info("handling text", slog.String("text", msg.Text), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Text {
	case "Посмотреть предстоящие события":
		return h.showEvents(ctx, msg)
	}
	return nil
}

// startMessage обработчик для команды /start
func (h *Handler) startMessage(ctx context.Context, msg *tgbotapi.Message) error {
	// Формируем ответ для пользователя
	answer := tgbotapi.NewMessage(msg.Chat.ID,
		"Привет! 👋\nЯ бот для отслеживания и регистрации на события.",
	)
	// После приветствия показываем основную клавиатуру
	answer.ReplyMarkup = keyboard.MainKeyboard()

	h.log.Info("saving user info", slog.Int("chat_id", int(msg.Chat.ID)), slog.String("username", msg.From.UserName), slog.Time("created_at", time.Now()))
	// Сохраняем информацию в базу данных
	if err := h.service.SaveUserInfo(ctx, msg.Chat.ID, msg.From.UserName); err != nil {
		return err
	}

	// Отправляем ответ пользователю
	_, err := h.sender.Send(answer)
	if err != nil {
		h.log.Error("error send answer", err.Error())
		return err
	}
	return nil
}

// showEvents обработчик для команды /getEvents
func (h *Handler) showEvents(ctx context.Context, msg *tgbotapi.Message) error {
	// Отправляем данные в сервисный слой, в случае ошибки - отправляем пользователю соответствующее сообщение
	events, err := h.service.GetEvents(ctx)
	if err != nil {
		errMsg := tgbotapi.NewMessage(msg.Chat.ID, "Произошла ошибка")
		_, err = h.sender.Send(errMsg)
		if err != nil {
			h.log.Error("error send answer", err.Error())
			return err
		}
		return err
	}

	if len(events) == 0 {
		noEventsMsg := tgbotapi.NewMessage(msg.Chat.ID, "К сожалению, событий не найдено")
		noEventsMsg.ReplyMarkup = keyboard.MainKeyboard()
		_, err = h.sender.Send(noEventsMsg)
		if err != nil {
			h.log.Error("error send answer", err.Error())
			return err
		}
		return nil
	}

	// Создаём "кнопки" с соответствующими данными
	var eventsButtons []keyboard.EventButton
	for _, e := range events {
		eventsButtons = append(eventsButtons, keyboard.EventButton{
			EventID: e.Id,
			Title:   e.Title,
		})
	}

	message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите событие, для просмотра детальной информации")
	message.ParseMode = tgbotapi.ModeMarkdown
	// Отправляем клавиатуру с событиями пользователю
	message.ReplyMarkup = keyboard.EventsKeyboard(eventsButtons)
	_, err = h.sender.Send(message)
	if err != nil {
		h.log.Error("error send answer", err.Error())
		return err
	}
	return nil
}

// TODO: Реализовать метод для /getEvent
