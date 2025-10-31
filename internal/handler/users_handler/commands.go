package usershandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/handler/builders"
	"github.com/Egor213/notifyBot/internal/service/srverrs"
	"github.com/Egor213/notifyBot/pkg/utils"
	"github.com/Egor213/notifyBot/pkg/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (h *UserHandler) handleStart() (string, entity.ReplyMarkup) {
	keyboard := builders.BuildInlineKeyboard(
		entity.InlineKeyboard{Name: "Зарегистрироваться", Command: "register"},
		entity.InlineKeyboard{Name: "Проверить email", Command: "get_email"},
	)
	return "👋 Привет! Этот бот поможет вам управлять уведомлениями.\n\nДля начала зарегистрируйтесь, используя команду:\n`/register your@email.com`", keyboard
}

func (h *UserHandler) handleStartRegister(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email := msg.CommandArguments()

	if err := validation.ValidateEmail(email); err != nil {
		switch err {
		case validation.ErrEmptyEmail:
			return "📭 Пожалуйста, укажите email в формате:\n`/register your@email.com`", nil
		case validation.ErrInvalidEmail:
			return "⚠️ Некорректный email. Укажите корректный адрес:\n`/register your@email.com`", nil
		}
	}

	_, err := h.UserService.GetEmail(ctx, msg.Chat.ID)

	if err == nil {
		return fmt.Sprintf("⚠️ Пользователь с TG ID %d уже зарегистрирован.", msg.Chat.ID), nil
	} else {
		switch {
		case errors.Is(err, srverrs.ErrUserNotFound):
			code := utils.GenerateCode()
			h.StateService.SetState(msg.Chat.ID, entity.StateAwaitingVerificationCode, map[any]any{
				msg.Chat.ID: email,
				email:       code,
			})
			// Почему то не работают коды доступа приложении
			go h.MailService.SendMessage(email, "NotifyBot Code", fmt.Sprintf("Code: %s", code))
			log.Infof("CODE: %s", code)
			return fmt.Sprintf(
				"📩 На почту %s отправлен код подтверждения.\n"+
					"Пожалуйста, введите его здесь, чтобы завершить регистрацию ✅",
				email,
			), nil
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "⚠️ Ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}
}

func (h *UserHandler) handlerVerifyEmail(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	userState := h.StateService.GetState(msg.Chat.ID)

	if userState.State != entity.StateAwaitingVerificationCode {
		return "⚠️ У вас нет активного процесса регистрации. Пожалуйста, используйте команду /register.", nil
	}

	inputCode := msg.Text
	emailAny := userState.Data[msg.Chat.ID]
	email, _ := emailAny.(string)

	codeAny := userState.Data[email]
	code, _ := codeAny.(string)
	if code != inputCode {
		return "❌ Код неверный. Попробуйте ещё раз.", nil
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

	h.StateService.ClearState(msg.Chat.ID)

	keyboard := builders.BuildInlineKeyboard(
		entity.InlineKeyboard{Name: "Посмотреть email", Command: "get_email"},
		entity.InlineKeyboard{Name: "Настройки уведомлений", Command: "view_settings"},
	)

	return fmt.Sprintf("✅ Регистрация прошла успешно!\nВаша почта: %s", user.Email), keyboard

}

func (h *UserHandler) handleGetEmail(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	email, err := h.UserService.GetEmail(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrUserNotFound):
			keyboard := builders.BuildInlineKeyboard(
				entity.InlineKeyboard{Name: "Зарегистрироваться", Command: "register"},
			)
			return "🙁 Вы ещё не зарегистрированы.", keyboard
		case errors.Is(err, srverrs.ErrUserCheckFailed):
			return "⚠️ Произошла ошибка при проверке пользователя. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	keyboard := builders.BuildInlineKeyboard(
		entity.InlineKeyboard{Name: "Изменить настройки уведомлений", Command: "view_settings"},
	)

	return fmt.Sprintf("📨 Ваш текущий email: %s", email), keyboard
}
