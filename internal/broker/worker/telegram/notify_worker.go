package telegramworker

import (
	"context"

	"github.com/Egor213/notifyBot/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/segmentio/kafka-go"
)

type NotifyWorker struct {
	TgBot *bot.Bot
}

func NewNotifyWorker(b *bot.Bot) *NotifyWorker {
	return &NotifyWorker{b}
}

func (w *NotifyWorker) ProcessMsg(ctx context.Context, msg kafka.Message) {
	logMsg := string(msg.Value)
	parcedMsg := ParceLogMsg(logMsg)
	tgLogMsg := CreateTgLogMsg(parcedMsg)
	tgMsg := tgbotapi.NewMessage(1573846092, tgLogMsg)
	tgMsg.ParseMode = tgbotapi.ModeMarkdown
	w.TgBot.Tg.Send(tgMsg)
}
