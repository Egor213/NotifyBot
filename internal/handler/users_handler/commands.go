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
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Зарегистрироваться", "register"),
			tgbotapi.NewInlineKeyboardButtonData("Проверить email", "get_email"),
		),
	)
	return "👋 Привет! Этот бот поможет вам управлять уведомлениями.\n\nДля начала зарегистрируйтесь, используя команду:\n`/register your@email.com`", keyboard
}

func (h *UserHandler) handleRegister(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email := msg.CommandArguments()

	if err := validation.ValidateEmail(email); err != nil {
		switch err {
		case validation.ErrEmptyEmail:
			return "📭 Пожалуйста, укажите email в формате:\n`/register your@email.com`", nil
		case validation.ErrInvalidEmail:
			return "⚠️ Некорректный email. Укажите корректный адрес:\n`/register your@email.com`", nil
		}
	}

	user, err := h.UserService.RegisterUser(ctx, msg.Chat.ID, email)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserAlreadyExists):
			return fmt.Sprintf("⚠️ Пользователь с TG ID %d уже зарегистрирован.", msg.Chat.ID), nil
		case errors.Is(err, srverrs.ErrUserCreateFailed):
			return "❌ Произошла ошибка при создании пользователя. Попробуйте позже.", nil
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "⚠️ Ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть email", "get_email"),
			tgbotapi.NewInlineKeyboardButtonData("Настройки уведомлений", "view_settings"),
		),
	)

	return fmt.Sprintf("✅ Регистрация прошла успешно!\nВаша почта: %s", user.Email), keyboard
}

func (h *UserHandler) handleGetEmail(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email, err := h.UserService.GetEmail(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserNotFound):
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Зарегистрироваться", "register"),
				),
			)
			return "🙁 Вы ещё не зарегистрированы.", keyboard
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "⚠️ Произошла ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить настройки уведомлений", "view_settings"),
		),
	)

	return fmt.Sprintf("📨 Ваш текущий email: %s", email), keyboard
}
