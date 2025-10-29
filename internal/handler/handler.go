package handler

import (
	notifyhandler "github.com/Egor213/notifyBot/internal/handler/notify_handler"
	usershandler "github.com/Egor213/notifyBot/internal/handler/users_handler"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	CanHandle(command string) bool
	HandleCommand(msg *tgbotapi.Message) string
}

type Handler struct {
	handlers []CommandHandler
}

func NewHandler(handlers ...CommandHandler) *Handler {
	return &Handler{handlers: handlers}
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) string {
	if !msg.IsCommand() {
		return "Пожалуйста, используйте команды, начиная с '/'."
	}

	command := msg.Command()

	for _, handler := range h.handlers {
		if handler.CanHandle(command) {
			return handler.HandleCommand(msg)
		}
	}

	return "Неизвестная команда. Используйте /start для справки."
}

func ConfigureHandler(services *service.Services) *Handler {
	userHandler := usershandler.NewUserHandler(services.User)
	notifyHandler := notifyhandler.NewNotificationHandler()
	return NewHandler(userHandler, notifyHandler)
}
