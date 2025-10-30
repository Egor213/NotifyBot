package telegramworker

import (
	"fmt"
	"strings"

	"github.com/Egor213/notifyBot/pkg/utils"
)

func ParceLogMsg(msg string) map[string]string {
	parts := strings.Split(msg, ",")

	data := make(map[string]string)

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			data[key] = strings.ToLower(val)
		}
	}
	return data
}

func CreateTgLogMsg(data map[string]string) string {
	levelEmoji := map[string]string{
		"info":  "ℹ️",
		"warn":  "⚠️",
		"error": "❌",
	}

	emoji := levelEmoji[data["level"]]

	tgMessage := fmt.Sprintf(
		"%s *Service:* `%s`\n%s *Level:* `%s`\n📝 *Message:*\n```\n%s\n```",
		emoji, utils.EscapeMarkdownV2(data["service"]),
		emoji, utils.EscapeMarkdownV2(data["level"]),
		utils.EscapeMarkdownV2(data["Message"]),
	)
	return tgMessage
}
