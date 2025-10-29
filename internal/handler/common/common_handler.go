package common

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseHandler struct {
	commands map[string]func(context.Context, *tgbotapi.Message) string
}

func (b *BaseHandler) CanHandle(command string) bool {
	_, ok := b.commands[command]
	return ok
}

func (b *BaseHandler) HandleCommand(ctx context.Context, msg *tgbotapi.Message) string {
	if handlerFunc, ok := b.commands[msg.Command()]; ok {
		return handlerFunc(ctx, msg)
	}
	return fmt.Sprintf("Команда %s не поддерживается.", msg.Command())
}

func (b *BaseHandler) RegisterCommand(cmd string, f func(context.Context, *tgbotapi.Message) string) {
	if b.commands == nil {
		b.commands = make(map[string]func(context.Context, *tgbotapi.Message) string)
	}
	b.commands[cmd] = f
}
