package handler

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	notifyhandler "github.com/Egor213/notifyBot/internal/handler/notify_handler"
	usershandler "github.com/Egor213/notifyBot/internal/handler/users_handler"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	CanHandle(command string) bool
	HandleCommand(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup)
	CanHandleCallback(cmd string) bool
	HandleCallback(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup)
}

type Handler struct {
	commandHandlers []CommandHandler
	ctx             context.Context
}

func NewHandler(ctx context.Context, handlers ...CommandHandler) *Handler {
	return &Handler{ctx: ctx, commandHandlers: handlers}
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	if !msg.IsCommand() {
		return "Пожалуйста, используйте команды, начиная с '/'.", nil
	}

	command := msg.Command()
	for _, handler := range h.commandHandlers {
		if handler.CanHandle(command) {
			return handler.HandleCommand(h.ctx, msg)
		}
	}

	return "Неизвестная команда. Используйте /start для справки.", nil
}

func (h *Handler) HandleCallback(cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
	data := cb.Data
	for _, handler := range h.commandHandlers {
		if handler.CanHandleCallback(data) {
			return handler.HandleCallback(h.ctx, cb)
		}
	}
	return "Неизвестная callback", nil
}

func ConfigureHandler(ctx context.Context, services *service.Services) *Handler {
	userHandler := usershandler.NewUserHandler(services.User)
	notifyHandler := notifyhandler.NewNotificationHandler(services.NotifySettings)

	return &Handler{
		ctx:             ctx,
		commandHandlers: []CommandHandler{userHandler, notifyHandler},
	}
}
