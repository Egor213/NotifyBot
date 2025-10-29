package notifyhandler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/service/srverrs"
	"github.com/Egor213/notifyBot/pkg/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *NotificationHandler) handleStatus(_ context.Context, msg *tgbotapi.Message) string {
	return "Бот работает стабильно 🟢"
}

func (h *NotificationHandler) handleViewSettings(ctx context.Context, msg *tgbotapi.Message) string {
	settings, err := h.NotifyServ.GetSettings(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrGetNotifySettings):
			return "Произошла ошибка при проверке настроек. Попробуйте позже."
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err)
		}
	}

	if len(settings) == 0 {
		return "У вас пока нет настроек уведомлений."
	}

	var lines []string
	for _, s := range settings {
		lines = append(lines, fmt.Sprintf("Сервис: %s | Уровень: %s", s.Service, s.Level))
	}

	return "Ваши настройки уведомлений:\n" + strings.Join(lines, "\n")
}

func (h *NotificationHandler) handleSetSettings(ctx context.Context, msg *tgbotapi.Message) string {
	args := strings.Fields(msg.CommandArguments())
	if len(args) == 0 {
		return "Usage: /set_notify_settings service1,service2 [level1,level2]"
	}

	services, err := validation.ParseServices(args[0])
	if err != nil {
		return "Please specify at least one service."
	}

	var levels []entity.LogLevel
	if len(args) > 1 {
		levels, err = validation.ParseLogLevels(args[1])
		if err != nil {
			return fmt.Sprintf("Invalid log levels: %v", err)
		}
	} else {
		levels, _ = validation.ParseLogLevels("")
	}

	if err := h.NotifyServ.SetSettings(ctx, msg.Chat.ID, services, levels); err != nil {
		return fmt.Sprintf("Failed to save settings: %v", err)
	}

	return "Notification settings successfully updated ✅"
}

func (h *NotificationHandler) handleRemoveSettings(ctx context.Context, msg *tgbotapi.Message) string {
	args := strings.Fields(msg.CommandArguments())
	if len(args) != 2 {
		return "Использование: /del_notify_settings service level"
	}

	level, ok := entity.ParseLogLevel(strings.ToUpper(args[1]))
	if !ok {
		return fmt.Sprintf("Неверный уровень логов: %s", args[1])
	}

	if err := h.NotifyServ.RemoveSettings(ctx, msg.Chat.ID, args[0], level); err != nil {
		switch {
		case errors.Is(err, srverrs.ErrNotifySettingNotFound):
			return "Указанная настройка не найдена."
		case errors.Is(err, srverrs.ErrRemoveNotifySettings):
			return "Произошла ошибка при удалении настройки. Попробуйте позже."
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err)
		}
	}

	return "Настройка успешно удалена ✅"
}
