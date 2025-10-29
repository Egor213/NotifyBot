package usershandler

import (
	"context"

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
	h.RegisterCommand("start", func(_ context.Context, _ *tgbotapi.Message) string {
		return h.handleStart()
	})
	h.RegisterCommand("register", h.handleRegister)
	h.RegisterCommand("get_email", h.handleGetEmail)
}
