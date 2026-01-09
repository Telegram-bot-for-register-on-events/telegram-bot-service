package keyboard

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// EventButton описывает информацию, содержащуюся в "кнопке"
type EventButton struct {
	EventID string
	Title   string
}

// MainKeyboard основная клавиатура "внизу экрана" для получения предстоящих событий
func MainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Посмотреть предстоящие события")))
}

// EventsKeyboard клавиатура "в сообщении" позволяет увидеть события
func EventsKeyboard(events []EventButton, numPage, pageSize, countEvents int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var navigationRow []tgbotapi.InlineKeyboardButton
	for _, event := range events {
		row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(event.Title, "event_"+event.EventID)}
		rows = append(rows, row)
	}

	if numPage > 0 {
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("Назад", "page_"+strconv.Itoa(numPage-1)))
	}

	if (numPage+1)*pageSize < countEvents {
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("Вперёд", "page_"+strconv.Itoa(numPage+1)))
	}

	if len(navigationRow) > 0 {
		rows = append(rows, navigationRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// EventDetailKeyboard клавиатура, которая показывает детальную информацию о событии, позволяет записаться на него
func EventDetailKeyboard(eventID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Зарегистрироваться", "register_"+eventID),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_events"),
		),
	)
}

// BackToSeeEvents возвращает к просмотру событий после регистрации
func BackToSeeEvents() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Продолжить просмотр событий", "back_to_events")))
}
