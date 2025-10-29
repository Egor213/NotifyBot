package entity

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type NotifySetting struct {
	ID      int64
	TgID    int64
	Service string
	Level   LogLevel
}
