package usershandler

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/handler/common"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandler struct {
	*common.BaseHandler
	UserService service.Users
}

func NewUserHandler(s service.Users) *UserHandler {
	h := &UserHandler{
		BaseHandler: &common.BaseHandler{},
		UserService: s,
	}
	h.registerCommands()
	return h
}

func (h *UserHandler) registerCommands() {
	h.RegisterCommand("start", func(_ context.Context, _ *tgbotapi.Message) (string, entity.ReplyMarkup) {
		return h.handleStart()
	})
	h.RegisterCommand("register", h.handleRegister)
	h.RegisterCommand("get_email", h.handleGetEmail)

	h.RegisterCallback("register", func(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
		return "✉️ Чтобы зарегистрироваться, введите:\n`/register your@email.com`", nil
	})

	h.RegisterCallback("get_email", func(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
		msg := cb.Message
		return h.handleGetEmail(ctx, msg)
	})

	h.RegisterCallback("view_settings", func(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
		return "⚙️ Чтобы посмотреть настройки, используйте команду `/notify_settings`", nil
	})
}
