package notifyhandler

import (
	"github.com/Egor213/notifyBot/internal/handler/common"
)

type NotificationHandler struct {
	*common.BaseHandler
}

func NewNotificationHandler() *NotificationHandler {
	h := &NotificationHandler{
		BaseHandler: &common.BaseHandler{},
	}
	h.registerCommands()
	return h
}

func (h *NotificationHandler) registerCommands() {
	h.RegisterCommand("notify", h.handleNotify)
	h.RegisterCommand("status", h.handleStatus)
}
