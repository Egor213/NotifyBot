package entity

import "strings"

func ParseLogLevel(s string) (LogLevel, bool) {
	switch strings.ToUpper(s) {
	case string(LogLevelInfo):
		return LogLevelInfo, true
	case string(LogLevelWarn):
		return LogLevelWarn, true
	case string(LogLevelError):
		return LogLevelError, true
	default:
		return "", false
	}
}
