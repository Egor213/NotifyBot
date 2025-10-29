package notifyhandler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationHandler struct {
	commands map[string]func(*tgbotapi.Message) string
}

func NewNotificationHandler() *NotificationHandler {
	h := &NotificationHandler{
		commands: make(map[string]func(*tgbotapi.Message) string),
	}
	h.registerCommands()
	return h
}

func (h *NotificationHandler) registerCommands() {
	h.commands["notify"] = h.handleNotify
	h.commands["status"] = h.handleStatus
}

func (h *NotificationHandler) CanHandle(command string) bool {
	_, ok := h.commands[command]
	return ok
}

func (h *NotificationHandler) HandleCommand(msg *tgbotapi.Message) string {
	if handlerFunc, ok := h.commands[msg.Command()]; ok {
		return handlerFunc(msg)
	}
	return fmt.Sprintf("Команда %s не поддерживается.", msg.Command())
}
