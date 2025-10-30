package telegramworker

import (
	"github.com/Egor213/notifyBot/internal/repository/repotypes"
)

func BuildChatIDFilter(data map[string]string) repotypes.ChatIDFilter {
	return repotypes.ChatIDFilter{
		Service: data["service"],
		Level:   data["level"],
	}
}
