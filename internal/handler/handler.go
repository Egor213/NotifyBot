package handler

import (
	"context"

	notifyhandler "github.com/Egor213/notifyBot/internal/handler/notify_handler"
	usershandler "github.com/Egor213/notifyBot/internal/handler/users_handler"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	CanHandle(command string) bool
	HandleCommand(ctx context.Context, msg *tgbotapi.Message) string
}

type Handler struct {
	handlers []CommandHandler
	ctx      context.Context
}

func NewHandler(ctx context.Context, handlers ...CommandHandler) *Handler {
	return &Handler{ctx: ctx, handlers: handlers}
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) string {
	if !msg.IsCommand() {
		return "Пожалуйста, используйте команды, начиная с '/'."
	}

	command := msg.Command()

	for _, handler := range h.handlers {
		if handler.CanHandle(command) {
			return handler.HandleCommand(h.ctx, msg)
		}
	}

	return "Неизвестная команда. Используйте /start для справки."
}

func ConfigureHandler(ctx context.Context, services *service.Services) *Handler {
	userHandler := usershandler.NewUserHandler(services.User)
	notifyHandler := notifyhandler.NewNotificationHandler(services.NotifySettings)
	return NewHandler(ctx, userHandler, notifyHandler)
}
