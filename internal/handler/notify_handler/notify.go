package notifyhandler

import (
	"github.com/Egor213/notifyBot/internal/handler/common"
	"github.com/Egor213/notifyBot/internal/service"
)

type NotificationHandler struct {
	*common.BaseHandler
	NotifyServ service.NotifySettings
}

func NewNotificationHandler(notityServ service.NotifySettings) *NotificationHandler {
	h := &NotificationHandler{
		BaseHandler: &common.BaseHandler{},
		NotifyServ:  notityServ,
	}
	h.registerCommands()
	return h
}

func (h *NotificationHandler) registerCommands() {
	h.RegisterCommand("status", h.handleStatus)
	h.RegisterCommand("notify_settings", h.handleViewSettings)
	h.RegisterCommand("set_notify_settings", h.handleSetSettings)
	h.RegisterCommand("del_notify_settings", h.handleRemoveSettings)
}
