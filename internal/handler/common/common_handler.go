package common

import (
	"context"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseHandler struct {
	commands  map[string]func(context.Context, *tgbotapi.Message) (string, entity.ReplyMarkup)
	callbacks map[string]func(context.Context, *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup)
}

func (b *BaseHandler) CanHandle(command string) bool {
	_, ok := b.commands[command]
	return ok
}

func (b *BaseHandler) HandleCommand(ctx context.Context, msg *tgbotapi.Message) (string, entity.ReplyMarkup) {
	if handlerFunc, ok := b.commands[msg.Command()]; ok {
		return handlerFunc(ctx, msg)
	}
	return fmt.Sprintf("Команда %s не поддерживается.", msg.Command()), nil
}

func (b *BaseHandler) RegisterCommand(cmd string, f func(context.Context, *tgbotapi.Message) (string, entity.ReplyMarkup)) {
	if b.commands == nil {
		b.commands = make(map[string]func(context.Context, *tgbotapi.Message) (string, entity.ReplyMarkup))
	}
	b.commands[cmd] = f
}

func (b *BaseHandler) RegisterCallback(key string, f func(context.Context, *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup)) {
	if b.callbacks == nil {
		b.callbacks = make(map[string]func(context.Context, *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup))
	}
	b.callbacks[key] = f
}

func (b *BaseHandler) CanHandleCallback(data string) bool {
	_, ok := b.callbacks[data]
	return ok
}

func (b *BaseHandler) HandleCallback(ctx context.Context, cb *tgbotapi.CallbackQuery) (string, entity.ReplyMarkup) {
	if handlerFunc, ok := b.callbacks[cb.Data]; ok {
		return handlerFunc(ctx, cb)
	}
	return fmt.Sprintf("Неизвестный callback: %s", cb.Data), nil
}
