package telegramworker

import (
	"context"

	"github.com/Egor213/notifyBot/internal/service"
	"github.com/Egor213/notifyBot/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type NotifyWorker struct {
	TgBot              *bot.Bot
	NotifySettingsServ service.NotifySettings
}

func NewNotifyWorker(b *bot.Bot, s service.NotifySettings) *NotifyWorker {
	return &NotifyWorker{b, s}
}

func (w *NotifyWorker) ProcessMsg(ctx context.Context, msg kafka.Message) {
	logMsg := string(msg.Value)
	parcedMsg := ParceLogMsg(logMsg)
	filter := BuildChatIDFilter(parcedMsg)
	chatIDs, err := w.NotifySettingsServ.GetChatIDsByFilters(ctx, filter)
	if err != nil {
		log.Errorf("NotifyWorker.ProcessMsg - w.NotifySettingsServ.GetChatIDsByFilters: %w", err)
		return
	}
	tgLogMsg := CreateTgLogMsg(parcedMsg)
	for _, chatID := range chatIDs {
		tgMsg := tgbotapi.NewMessage(chatID, tgLogMsg)
		tgMsg.ParseMode = tgbotapi.ModeMarkdown
		w.TgBot.Tg.Send(tgMsg)
	}
}
