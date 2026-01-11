package bot

import (
	"log/slog"
	"time"

	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/bot/handlers"
	"github.com/Telegram-bot-for-register-on-events/telegram-bot-service/internal/service"
	tele "gopkg.in/telebot.v3"
)

// Bot описывает телеграм-бота
type Bot struct {
	log     *slog.Logger
	bot     *tele.Bot
	handler *handlers.Handler
}

// NewBot конструктор для Bot
func NewBot(log *slog.Logger, token string, service *service.Service) (*Bot, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	h := handlers.NewHandler(log, service, b)

	return &Bot{
		log:     log,
		bot:     b,
		handler: h,
	}, nil
}

// MustStart запускает бота, в прослойке с помощью recover отлавливает паники, регистрирует обработчики
func (b *Bot) MustStart() {
	b.bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			defer func() {
				if r := recover(); r != nil {
					b.log.Error("panic recovered", slog.Any("panic", r))
				}
			}()
			return next(c)
		}
	})

	b.handler.RegisterHandlers(b.bot)
	b.bot.Start()
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.bot.Stop()
}
