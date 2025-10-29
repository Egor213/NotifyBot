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

func (h *NotificationHandler) handleStatus(_ context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	return "Бот работает стабильно 🟢", nil
}

func (h *NotificationHandler) handleViewSettings(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {

	settings, err := h.NotifyServ.GetSettings(ctx, msg.Chat.ID)
	if err != nil {
		switch {
		case errors.Is(err, srverrs.ErrGetNotifySettings):
			return "Произошла ошибка при получении настроек уведомлений. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	if len(settings) == 0 {
		return "У вас пока нет настроек уведомлений.", nil
	}

	serviceMap := make(map[string][]string)
	for _, s := range settings {
		serviceMap[s.Service] = append(serviceMap[s.Service], string(s.Level))
	}

	var lines []string
	for svc, levels := range serviceMap {
		lines = append(lines, fmt.Sprintf("Сервис: %s | Уровни: %s", svc, strings.Join(levels, ", ")))
	}

	return "Ваши настройки уведомлений:\n" + strings.Join(lines, "\n"), nil
}

func (h *NotificationHandler) handleSetSettings(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) == 0 {
		return "Используйте: /set_notify_settings service1,service2 [level1,level2]", nil
	}

	services, err := validation.ParseServices(args[0])
	if err != nil {
		return "Необходимо вписать хотя бы 1 сервис.", nil
	}

	var levels []entity.LogLevel
	if len(args) > 1 {
		levels, err = validation.ParseLogLevels(args[1])
		if err != nil {
			return fmt.Sprintf("Неверный уровень логирования: %v", err), nil
		}
	} else {
		levels, _ = validation.ParseLogLevels("")
	}

	if err := h.NotifyServ.SetSettings(ctx, msg.Chat.ID, services, levels); err != nil {
		return fmt.Sprintf("Не удалось сохранить настройки: %v", err), nil
	}

	return "Настройки оповещений установлены ✅", nil
}

func (h *NotificationHandler) handleRemoveSettings(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) != 2 {
		return "Использование: /del_notify_settings service level", nil
	}

	level, ok := entity.ParseLogLevel(strings.ToUpper(args[1]))
	if !ok {
		return fmt.Sprintf("Неверный уровень логов: %s", args[1]), nil
	}

	if err := h.NotifyServ.RemoveSettings(ctx, msg.Chat.ID, args[0], level); err != nil {
		switch {
		case errors.Is(err, srverrs.ErrNotifySettingNotFound):
			return "Указанная настройка не найдена.", nil
		case errors.Is(err, srverrs.ErrRemoveNotifySettings):
			return "Произошла ошибка при удалении настройки. Попробуйте позже.", nil
		default:
			return fmt.Sprintf("Неизвестная ошибка: %v", err), nil
		}
	}

	return "Настройка успешно удалена ✅", nil
}
