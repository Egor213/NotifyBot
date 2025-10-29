package usershandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/service/srverrs"
	"github.com/Egor213/notifyBot/pkg/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *UserHandler) handleStart() string {
	return "Привет! Чтобы зарегистрироваться, отправьте: /register your@email.com"
}

func (h *UserHandler) handleRegister(ctx context.Context, msg *tgbotapi.Message) string {
	email := msg.CommandArguments()

	if err := validation.ValidateEmail(email); err != nil {
		switch err {
		case validation.ErrEmptyEmail:
			return "Пожалуйста, укажите email: /register your@email.com"
		case validation.ErrInvalidEmail:
			return "Некорректный email. Укажите корректный адрес: /register your@email.com"
		}
	}

	user, err := h.UserService.RegisterUser(ctx, msg.Chat.ID, email)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserAlreadyExists):
			return fmt.Sprintf("Пользователь с TG ID %d уже зарегистрирован.", msg.Chat.ID)
		case errors.Is(err, srverrs.ErrUserCreateFailed):
			return "Произошла ошибка при создании пользователя. Попробуйте позже."
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "Произошла ошибка при проверке пользователя. Попробуйте позже."
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err)
		}
	}

	return fmt.Sprintf("Вы успешно зарегистрированы! Ваша почта: %s", user.Email)
}

func (h *UserHandler) handleGetEmail(ctx context.Context, msg *tgbotapi.Message) string {
	email, err := h.UserService.GetEmail(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserNotFound):
			return "Вы ещё не зарегистрированы. Используйте /register your@email.com"
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "Произошла ошибка при проверке пользователя. Попробуйте позже."
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err)
		}
	}

	return fmt.Sprintf("Ваш текущий email: %s", email)
}
