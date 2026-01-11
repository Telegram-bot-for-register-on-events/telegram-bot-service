package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/bot/keyboard"
	tele "gopkg.in/telebot.v3"
)

// Service описывает методы для взаимодействия с сервисным слоем
type Service interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
	SaveUserInfo(ctx context.Context, chatID int64, username string) error
}

// Handler описывает слой обработчиков
type Handler struct {
	log     *slog.Logger
	service Service
}

// NewHandler конструктор для Handler
func NewHandler(log *slog.Logger, service Service, _ *tele.Bot) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

// RegisterHandlers регистрирует обработчики для клавиатур и комманд
func (h *Handler) RegisterHandlers(b *tele.Bot) {
	b.Handle("/start", h.startMessage)
	b.Handle(tele.OnText, h.handleText)
	b.Handle(tele.OnCallback, h.handleCallback)
}

// startMessage обработчик для команды /start
func (h *Handler) startMessage(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	chatID := c.Chat().ID
	username := c.Sender().Username

	h.log.Info("saving user info", slog.Int64("chat_id", chatID), slog.String("username", username))
	if err := h.service.SaveUserInfo(ctx, chatID, username); err != nil {
		h.log.Error("failed to save user", slog.String("error", err.Error()))
	}

	return c.Send(
		"Привет! 👋\nЯ бот для отслеживания и регистрации на события.",
		keyboard.MainKeyboard(),
	)
}

// handleText обработчик для текстовых сообщений
func (h *Handler) handleText(c tele.Context) error {
	if c.Text() == "Посмотреть предстоящие события" {
		return h.showEvents(c, 0)
	}
	return nil
}

// showEvents показывает список событий
func (h *Handler) showEvents(c tele.Context, pageNum int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	events, err := h.service.GetEvents(ctx)
	if err != nil {
		return c.Send("Ошибка при получении событий")
	}

	h.log.Info("events from service", slog.Int("count", len(events)))

	if len(events) == 0 {
		return c.Send("Событий не найдено")
	}

	pageSize := 5
	totalEvents := len(events)
	start := pageNum * pageSize

	if start >= totalEvents {
		start = 0
		pageNum = 0
	}

	end := start + pageSize
	if end > totalEvents {
		end = totalEvents
	}

	var buttons []keyboard.EventButton
	for i := start; i < end; i++ {
		e := events[i]
		buttons = append(buttons, keyboard.EventButton{
			EventID: e.Id,
			Title:   e.Title,
		})
	}

	markup := keyboard.EventsKeyboard(buttons, pageNum, pageSize, totalEvents)

	if c.Callback() != nil {
		return c.Edit(
			"Выберите событие:",
			&tele.SendOptions{
				ParseMode:   tele.ModeMarkdown,
				ReplyMarkup: markup,
			},
		)
	}

	return c.Send(
		"Выберите событие:",
		&tele.SendOptions{
			ParseMode:   tele.ModeMarkdown,
			ReplyMarkup: markup,
		},
	)
}

// formatEventInfo форматирует строку с деталями информации
func formatEventInfo(e *pb.Event) string {
	t := e.StartsAt.AsTime().Format("02.01.2006 15:04")
	return fmt.Sprintf("*%s*\n\n%s\n\n*Начало:* %s",
		e.GetTitle(), e.GetDescription(), t)
}

// showEventDetails показывает детали события
func (h *Handler) showEventDetails(c tele.Context, eventID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.log.Info("showing event details", slog.String("event_id", eventID))

	event, err := h.service.GetEvent(ctx, eventID)
	if err != nil || event == nil {
		return h.showEvents(c, 0)
	}

	text := formatEventInfo(event)
	markup := keyboard.EventDetailKeyboard(eventID)

	return c.Edit(
		text,
		&tele.SendOptions{
			ParseMode:   tele.ModeMarkdown,
			ReplyMarkup: markup,
		},
	)
}

// backToEvents возвращает назад к просмотру событий
func (h *Handler) backToEvents(c tele.Context) error {
	return h.showEvents(c, 0)
}

// register регистрирует пользователя на событие
func (h *Handler) register(c tele.Context, eventID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := c.Sender()

	success, err := h.service.RegisterUser(ctx, eventID, c.Chat().ID, user.Username)
	if err != nil {
		return c.Edit(
			"Произошла ошибка.",
			&tele.SendOptions{
				ParseMode:   tele.ModeMarkdown,
				ReplyMarkup: keyboard.EventDetailKeyboard(eventID),
			},
		)
	}

	if success {
		return c.Edit(
			"Вы успешно зарегистрированы на это событие!",
			&tele.SendOptions{
				ParseMode:   tele.ModeMarkdown,
				ReplyMarkup: keyboard.BackToSeeEvents(),
			},
		)
	}

	return c.Edit(
		"Не удалось зарегистрироваться. Возможно, вы уже зарегистрированы на это событие.",
		&tele.SendOptions{
			ParseMode:   tele.ModeMarkdown,
			ReplyMarkup: keyboard.BackToSeeEvents(),
		},
	)
}

// handleCallback обработчик callback'ов
func (h *Handler) handleCallback(c tele.Context) error {
	callback := c.Callback()

	h.log.Info("callback received", slog.String("data", callback.Data), slog.Int64("chat_id", c.Chat().ID))

	if err := c.Respond(); err != nil {
		h.log.Error("failed to respond to callback", slog.String("error", err.Error()))
	}

	parts := strings.SplitN(callback.Data, ":", 2)
	if len(parts) < 2 {
		h.log.Error("invalid callback format", slog.String("data", callback.Data))
		return h.showEvents(c, 0)
	}

	action := parts[0]
	data := parts[1]

	h.log.Info("parsed callback", slog.String("action", action), slog.String("data", data))

	switch action {
	case "event":
		return h.showEventDetails(c, data)

	case "page":
		page, err := strconv.Atoi(data)
		if err != nil {
			h.log.Error("invalid page number", slog.String("data", data))
			return h.showEvents(c, 0)
		}
		return h.showEvents(c, page)

	case "back":
		return h.backToEvents(c)

	case "register":
		return h.register(c, data)

	default:
		h.log.Warn("unknown callback action", slog.String("action", action))
		return h.showEvents(c, 0)
	}
}
