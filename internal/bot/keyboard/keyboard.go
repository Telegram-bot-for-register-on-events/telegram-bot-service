package keyboard

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type EventButton struct {
	EventID string
	Title   string
}

func MainKeyboard() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{ResizeKeyboard: true}
	kb.Reply(kb.Row(tele.Btn{Text: "Посмотреть предстоящие события"}))
	return kb
}

func EventsKeyboard(events []EventButton, numPage, pageSize, countEvents int) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	var rows [][]tele.InlineButton

	for _, e := range events {
		btn := tele.InlineButton{
			Text: e.Title,
			Data: "event:" + e.EventID,
		}
		rows = append(rows, []tele.InlineButton{btn})
	}

	if numPage > 0 || (numPage+1)*pageSize < countEvents {
		var navRow []tele.InlineButton

		if numPage > 0 {
			navRow = append(navRow, tele.InlineButton{
				Text: "Назад",
				Data: "page:" + strconv.Itoa(numPage-1),
			})
		}

		if (numPage+1)*pageSize < countEvents {
			navRow = append(navRow, tele.InlineButton{
				Text: "Вперёд",
				Data: "page:" + strconv.Itoa(numPage+1),
			})
		}

		if len(navRow) > 0 {
			rows = append(rows, navRow)
		}
	}

	kb.InlineKeyboard = rows
	return kb
}

func EventDetailKeyboard(eventID string) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	kb.InlineKeyboard = [][]tele.InlineButton{
		{
			{Text: "Зарегистрироваться", Data: "register:" + eventID},
			{Text: "Назад к событиям", Data: "back:"},
		},
	}

	return kb
}

func BackToSeeEvents() *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}

	kb.InlineKeyboard = [][]tele.InlineButton{
		{
			{Text: "Продолжить просмотр событий", Data: "back:"},
		},
	}

	return kb
}
