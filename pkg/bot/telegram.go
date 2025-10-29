// pkg/bot/telegram.go
package bot

import (
	"log"

	"github.com/Egor213/notifyBot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tg      *tgbotapi.BotAPI
	handler *handler.Handler
	workers int
}

func NewBot(token string, h *handler.Handler, workers int, debug bool) *Bot {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	botAPI.Debug = debug

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать работу с ботом"},
		{Command: "register", Description: "Зарегистрировать email"},
		{Command: "notify", Description: "Отправить уведомление"},
		{Command: "status", Description: "Проверить статус бота"},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := botAPI.Request(cfg); err != nil {
		log.Panic(err)
	}

	return &Bot{
		tg:      botAPI,
		handler: h,
		workers: workers,
	}
}

func (b *Bot) Start(bot_timeout int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = bot_timeout

	updates := b.tg.GetUpdatesChan(u)
	msgChan := make(chan *tgbotapi.Message, 100)

	for i := 0; i < b.workers; i++ {
		go sendUpdMes(msgChan, b)
	}
	for update := range updates {
		if update.Message != nil {
			msgChan <- update.Message
		}
	}
}

func sendUpdMes(msgChan <-chan *tgbotapi.Message, b *Bot) {
	for msg := range msgChan {
		response := b.handler.HandleMessage(msg)
		msgToSend := tgbotapi.NewMessage(msg.Chat.ID, response)
		b.tg.Send(msgToSend)
	}
}
