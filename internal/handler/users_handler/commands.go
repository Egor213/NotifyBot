package usershandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/service/srverrs"
	"github.com/Egor213/notifyBot/pkg/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *UserHandler) handleStart() (string, entity.ReplyMarkup) {
	return "Привет! Чтобы зарегистрироваться, отправьте: /register your@email.com", nil
}

func (h *UserHandler) handleRegister(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email := msg.CommandArguments()

	if err := validation.ValidateEmail(email); err != nil {
		switch err {
		case validation.ErrEmptyEmail:
			return "Пожалуйста, укажите email: /register your@email.com", nil
		case validation.ErrInvalidEmail:
			return "Некорректный email. Укажите корректный адрес: /register your@email.com", nil
		}
	}

	user, err := h.UserService.RegisterUser(ctx, msg.Chat.ID, email)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserAlreadyExists):
			return fmt.Sprintf("Пользователь с TG ID %d уже зарегистрирован.", msg.Chat.ID), nil
		case errors.Is(err, srverrs.ErrUserCreateFailed):
			return "Произошла ошибка при создании пользователя. Попробуйте позже.", nil
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "Произошла ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	return fmt.Sprintf("Вы успешно зарегистрированы! Ваша почта: %s", user.Email), nil
}

func (h *UserHandler) handleGetEmail(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email, err := h.UserService.GetEmail(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserNotFound):
			return "Вы ещё не зарегистрированы. Используйте /register your@email.com", nil
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "Произошла ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	return fmt.Sprintf("Ваш текущий email: %s", email), nil
}
