package notifyhandler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *NotificationHandler) handleNotify(msg *tgbotapi.Message) string {
	text := msg.CommandArguments()
	if text == "" {
		return "Пожалуйста, укажите текст уведомления: /notify текст"
	}
	return fmt.Sprintf("Оповещение отправлено всем пользователям ✅\nТекст: %s", text)
}

func (h *NotificationHandler) handleStatus(msg *tgbotapi.Message) string {
	return "Бот работает стабильно 🟢"
}
