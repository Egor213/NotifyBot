package validation

import (
	"fmt"
	"strings"

	"github.com/Egor213/notifyBot/internal/entity"
)

func ParseServices(s string) ([]string, error) {
	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("at least one service must be specified")
	}
	return res, nil
}

func ParseLogLevels(levels string) ([]entity.LogLevel, error) {
	if levels == "" {
		return []entity.LogLevel{
			entity.LogLevelInfo,
			entity.LogLevelWarn,
			entity.LogLevelError,
		}, nil
	}

	parts := strings.Split(levels, ",")
	res := make([]entity.LogLevel, 0, len(parts))
	for _, l := range parts {
		l = strings.TrimSpace(strings.ToUpper(l))
		lvl, ok := entity.ParseLogLevel(l)
		if !ok {
			return nil, fmt.Errorf("invalid log level: %s", l)
		}
		res = append(res, lvl)
	}

	return res, nil
}
