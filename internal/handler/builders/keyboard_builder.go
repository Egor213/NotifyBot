package builders

import (
	"github.com/Egor213/notifyBot/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildInlineKeyboard(buttons ...entity.InlineKeyboard) tgbotapi.InlineKeyboardMarkup {
	const maxPerRow = 4
	rows := [][]tgbotapi.InlineKeyboardButton{}

	for i := 0; i < len(buttons); i += maxPerRow {
		end := i + maxPerRow
		if end > len(buttons) {
			end = len(buttons)
		}

		row := []tgbotapi.InlineKeyboardButton{}
		for _, btn := range buttons[i:end] {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(btn.Name, btn.Command))
		}
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
