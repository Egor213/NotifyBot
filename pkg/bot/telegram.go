// pkg/bot/telegram.go
package bot

import (
	log "github.com/sirupsen/logrus"

	"github.com/Egor213/notifyBot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getCommands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать работу с ботом"},
		{Command: "register", Description: "Зарегистрировать email"},
		{Command: "status", Description: "Проверить статус бота"},
		{Command: "get_email", Description: "Посмотреть свою почту"},
		{Command: "notify_settings", Description: "Настройки оповещений"},
		{Command: "set_notify_settings", Description: "Устновить настройки оповещений"},
		{Command: "del_notify_settings", Description: "Удалить настройки оповещений"},
	}
}

type Bot struct {
	Tg      *tgbotapi.BotAPI
	handler *handler.Handler
	workers int
}

func NewBot(token string, h *handler.Handler, workers int, debug bool) *Bot {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	botAPI.Debug = debug

	commands := getCommands()

	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := botAPI.Request(cfg); err != nil {
		log.Panic(err)
	}

	return &Bot{
		Tg:      botAPI,
		handler: h,
		workers: workers,
	}
}

func (b *Bot) Start(bot_timeout int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = bot_timeout

	updates := b.Tg.GetUpdatesChan(u)
	msgChan := make(chan *tgbotapi.Message, 100)
	cbChan := make(chan *tgbotapi.CallbackQuery, 100)

	for i := 0; i < b.workers; i++ {
		go sendUpdMes(msgChan, b)
		go handleCallback(cbChan, b)
	}

	for update := range updates {
		if update.Message != nil {
			msgChan <- update.Message
		}
		if update.CallbackQuery != nil {
			cbChan <- update.CallbackQuery
		}
	}
}

func sendUpdMes(msgChan <-chan *tgbotapi.Message, b *Bot) {
	for msg := range msgChan {
		response, keyboard := b.handler.HandleMessage(msg)
		msgToSend := tgbotapi.NewMessage(msg.Chat.ID, response)
		msgToSend.ReplyMarkup = keyboard
		msgToSend.ParseMode = tgbotapi.ModeMarkdown
		b.Tg.Send(msgToSend)
	}
}

func handleCallback(cbChan <-chan *tgbotapi.CallbackQuery, b *Bot) {
	for cb := range cbChan {
		b.Tg.Request(tgbotapi.NewCallback(cb.ID, ""))

		chatID := cb.Message.Chat.ID
		response, keyboard := b.handler.HandleCallback(cb)
		msgToSend := tgbotapi.NewMessage(chatID, response)
		msgToSend.ParseMode = tgbotapi.ModeMarkdown
		msgToSend.ReplyMarkup = keyboard
		b.Tg.Send(msgToSend)
	}
}
