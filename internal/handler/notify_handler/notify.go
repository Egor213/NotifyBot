package notifyhandler

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/handler/common"
	"github.com/Egor213/notifyBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	h.RegisterCallback("view_settings", h.handleViewSettingsCallback)
	h.RegisterCallback("remove_settings", h.handleRemoveSettingsCallback)
	h.RegisterCallback("set_settings", func(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
		return "🛠 Чтобы добавить настройки, используйте:\n`/set_notify_settings service1,service2 level1,level2`", nil
	})
}
