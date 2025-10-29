package usershandler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *UserHandler) handleStart() string {
	return "Привет! Чтобы зарегистрироваться, отправьте: /register your@email.com"
}

func (h *UserHandler) handleRegister(msg *tgbotapi.Message) string {
	args := msg.CommandArguments()
	if args == "" {
		return "Пожалуйста, укажите email: /register your@email.com"
	}

	user, err := h.service.RegisterUser(msg.Chat.ID, args)
	if err != nil {
		return fmt.Sprintf("Ошибка регистрации: %v", err)
	}

	return fmt.Sprintf("Вы зарегистрированы! Ваша почта: %s", user.Email)
}
