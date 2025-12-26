package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// EventButton описывает информацию, содержащуюся в "кнопке"
type EventButton struct {
	EventID string
	Title   string
}

// mainKeyboard основная клавиатура "внизу экрана" для получения предстоящих событий
func mainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Посмотреть предстоящие события")))
}

// eventsKeyboard клавиатура "в сообщении" позволяет увидеть события
func eventsKeyboard(events []EventButton) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, event := range events {
		row := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(event.EventID, event.Title)}
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// eventDetailKeyboard клавиатура, которая показывает детальную информацию о событии, позволяет записаться на него
func eventDetailKeyboard(eventID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Записаться", "event_"+eventID),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_events"),
		),
	)
}

// replyKeyboardRemove позволяет убрать inline клавиатуру (внизу экрана)
func replyKeyboardRemove() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}
