package handler

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	notifyhandler "github.com/Egor213/notifyBot/internal/handler/notify_handler"
	usershandler "github.com/Egor213/notifyBot/internal/handler/users_handler"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Как будто вот эти все хендлеры надо было сразу делать, чтобы они возращали сконструйрованное сообщение для бота, а не просто текст и клавиатуру
// В следующий раз буду знать)

type CommandHandler interface {
	CanHandle(command string) bool
	CanHandleState(state entity.StateType) bool
	CanHandleCallback(cmd string) bool
	HandleCommand(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup)
	HandleCallback(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup)
	HandleState(ctx context.Context, msg *tgbotapi.Message, state entity.StateType) (string, entity.ReplyMarkup)
}

type Handler struct {
	commandHandlers []CommandHandler
	ctx             context.Context
	StateServ       service.State
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
	// Тут можно сразу services передавать а не перебирать руками, все равно указатель
	userHandler := usershandler.NewUserHandler(services.User, services.State, services.MailSender)
	notifyHandler := notifyhandler.NewNotificationHandler(services.NotifySettings)

	return &Handler{
		ctx:             ctx,
		commandHandlers: []CommandHandler{userHandler, notifyHandler},
		StateServ:       services.State,
	}
}

func (h *Handler) HandleNonCommandMessage(msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	state := h.StateServ.GetState(msg.Chat.ID).State
	for _, handler := range h.commandHandlers {
		if handler.CanHandleState(state) {
			return handler.HandleState(h.ctx, msg, state)
		}
	}
	h.StateServ.ClearState(msg.Chat.ID)
	return "Не понимаю, что вы имеете в виду. Используйте /start.", nil
}
