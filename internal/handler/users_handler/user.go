package usershandler

import (
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandler struct {
	service  service.Users
	commands map[string]func(*tgbotapi.Message) string
}

func NewUserHandler(s service.Users) *UserHandler {
	h := &UserHandler{
		service:  s,
		commands: make(map[string]func(*tgbotapi.Message) string),
	}
	h.registerCommands()
	return h
}

func (h *UserHandler) registerCommands() {
	h.commands["start"] = func(msg *tgbotapi.Message) string { return h.handleStart() }
	h.commands["register"] = h.handleRegister
}

func (h *UserHandler) CanHandle(command string) bool {
	_, ok := h.commands[command]
	return ok
}

func (h *UserHandler) HandleCommand(msg *tgbotapi.Message) string {
	if handlerFunc, ok := h.commands[msg.Command()]; ok {
		return handlerFunc(msg)
	}
	return "Неизвестная команда пользователя."
}
