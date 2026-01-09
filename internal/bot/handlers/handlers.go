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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service описывает методы для взаимодействия с сервисным слоем
type Service interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
	SaveUserInfo(ctx context.Context, chatID int64, username string) error
}

// Sender описывает методы для отправки сообщений пользователю и telegram'у
type Sender interface {
	// Send для отправки сообщений пользователю
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	// Request для отправки callback'ов
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
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
	if update.CallbackQuery != nil {
		return h.handleCallbackQuery(ctx, update.CallbackQuery)
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

// handleCallbackQuery обработчик нажатий на inline-кнопки
func (h *Handler) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	h.log.Info("handling callback query", slog.String("data", callback.Data), slog.Int64("chat_id", callback.Message.Chat.ID))

	if err := h.answerCallback(callback.ID); err != nil {
		h.log.Error("failed to answer callback", slog.String("error", err.Error()))
	}
	
	dataSplit := strings.Split(callback.Data, "_")
	switch dataSplit[0] {
	case "event":
		return h.showEventDetails(ctx, callback)
	case "back":
		return h.backToEvents(ctx, callback)
	case "register":
		return h.register(ctx, callback)
	case "page":
		numPage, _ := strconv.Atoi(dataSplit[1])
		return h.showEventsPage(ctx, callback.Message.Chat.ID, callback.Message.MessageID, numPage)
	}
	return nil
}

// answerCallback отвечает telegram, что callback получен
func (h *Handler) answerCallback(callbackID string) error {
	answer := tgbotapi.NewCallback(callbackID, "")
	_, err := h.sender.Request(answer)
	if err != nil {
		h.log.Error("error answer callback", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// handleCommand метод для обработки команд
func (h *Handler) handleCommand(ctx context.Context, msg *tgbotapi.Message) error {
	h.log.Info("handling command", slog.String("command", msg.Command()), slog.Int("chat_id", int(msg.Chat.ID)))
	switch msg.Command() {
	case "start":
		return h.startMessage(ctx, msg)
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
		h.log.Error("error answer on command", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// showEvents обработчик для команды /getEvents
func (h *Handler) showEvents(ctx context.Context, msg *tgbotapi.Message) error {
	eventsButtons, countEvents, err := h.gettingEventsForPage(ctx, msg.Chat.ID, 0)
	if err != nil {
		return err
	}

	message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите событие, для просмотра детальной информации")
	message.ReplyMarkup = keyboard.EventsKeyboard(eventsButtons, 0, 5, countEvents)
	message.ParseMode = tgbotapi.ModeMarkdown

	_, err = h.sender.Send(message)
	if err != nil {
		h.log.Error("error answer on command", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// showEventsPage отображает конкретную "страницу" со списком событий
func (h *Handler) showEventsPage(ctx context.Context, chatID int64, messageID int, numPage int) error {
	eventsButtons, countEvents, err := h.gettingEventsForPage(ctx, chatID, numPage)
	if err != nil {
		return err
	}

	message := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "Выберите событие, для просмотра детальной информации", keyboard.EventsKeyboard(eventsButtons, numPage, 5, countEvents))
	message.ParseMode = tgbotapi.ModeMarkdown

	_, err = h.sender.Send(message)
	if err != nil {
		h.log.Error("error answer callback", slog.String("error", err.Error()))
		return err
	}
	return nil

}

// gettingEventsForPage метод для получения событий и формирования рядов с ними
func (h *Handler) gettingEventsForPage(ctx context.Context, chatID int64, numPage int) ([]keyboard.EventButton, int, error) {
	// Отправляем данные в сервисный слой, в случае ошибки - отправляем пользователю соответствующее сообщение
	events, err := h.service.GetEvents(ctx)
	if err != nil {
		errMsg := tgbotapi.NewMessage(chatID, "Произошла ошибка")
		_, err = h.sender.Send(errMsg)
		if err != nil {
			h.log.Error("error send answer about error", slog.String("error", err.Error()))
			return nil, 0, err
		}
		return nil, 0, err
	}

	if len(events) == 0 {
		noEventsMsg := tgbotapi.NewMessage(chatID, "К сожалению, событий не найдено")
		noEventsMsg.ReplyMarkup = keyboard.MainKeyboard()
		_, err = h.sender.Send(noEventsMsg)
		if err != nil {
			h.log.Error("error send answer about error", slog.String("error", err.Error()))
			return nil, 0, err
		}
	}

	countEvents := len(events)
	start := numPage * 5
	end := start + 5

	if start >= countEvents {
		start = 0
	}
	if end > countEvents {
		end = countEvents
	}

	pageEvents := events[start:end]

	// Создаём "кнопки" с соответствующими данными
	var eventsButtons []keyboard.EventButton
	for _, e := range pageEvents {
		eventsButtons = append(eventsButtons, keyboard.EventButton{
			EventID: e.Id,
			Title:   e.Title,
		})
	}

	return eventsButtons, countEvents, nil
}

// showEventDetails показывает детали события
func (h *Handler) showEventDetails(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	dataSplit := strings.Split(callback.Data, "_")
	e, err := h.service.GetEvent(ctx, dataSplit[1])
	if err != nil {
		errEventsDetails := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка получения деталей события")
		errEventsDetails.ReplyMarkup = keyboard.MainKeyboard()
		_, err = h.sender.Send(errEventsDetails)
		if err != nil {
			h.log.Error("error send answer about error", slog.String("error", err.Error()))
			return err
		}
	}

	// Форматируем информацию
	eventInfo := formatEventInfo(e)

	// Заменяем прошлое сообщение
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(callback.Message.Chat.ID, callback.Message.MessageID, eventInfo, keyboard.EventDetailKeyboard(dataSplit[1]))
	editMsg.ParseMode = tgbotapi.ModeMarkdown

	_, err = h.sender.Send(editMsg)
	return err
}

// backToEvents обработчик для Inline-кнопки "назад" при просмотре деталей события
func (h *Handler) backToEvents(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	eventsButtons, countEvents, err := h.gettingEventsForPage(ctx, callback.Message.Chat.ID, 0)
	if err != nil {
		return err
	}

	// Заменяем прошлое сообщение
	message := tgbotapi.NewEditMessageTextAndMarkup(callback.Message.Chat.ID, callback.Message.MessageID, "Выберите событие, для просмотра детальной информации", keyboard.EventsKeyboard(eventsButtons, 0, 5, countEvents))
	_, err = h.sender.Send(message)
	if err != nil {
		h.log.Error("error answer callback", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// register регистрирует пользователя на конкретное событие
func (h *Handler) register(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	dataSplit := strings.Split(callback.Data, "_")
	result, err := h.service.RegisterUser(ctx, dataSplit[1], callback.Message.Chat.ID, callback.Message.Chat.UserName)
	if err != nil {
		errMsg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка. Попробуйте ещё раз")
		_, err = h.sender.Send(errMsg)
		if err != nil {
			h.log.Error("error answer about error", slog.String("error", err.Error()))
			return err
		}
		return err
	}

	if result {
		message := tgbotapi.NewEditMessageTextAndMarkup(callback.Message.Chat.ID, callback.Message.MessageID, "Вы успешно зарегистрированы!", keyboard.BackToSeeEvents())
		_, err = h.sender.Send(message)
		if err != nil {
			h.log.Error("error answer on callback", slog.String("error", err.Error()))
			return err
		}
	} else {
		message := tgbotapi.NewEditMessageTextAndMarkup(callback.Message.Chat.ID, callback.Message.MessageID, "Произошла ошибка", keyboard.BackToSeeEvents())
		_, err = h.sender.Send(message)
		if err != nil {
			h.log.Error("error answer about error", slog.String("error", err.Error()))
			return err
		}
	}
	return nil
}

// formatEventInfo форматирует информацию о событии
func formatEventInfo(e *pb.Event) string {
	t := e.StartsAt.AsTime().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("Название: %s\n Описание: %s\n Начало: %s", e.Title, e.Description, t)
}
