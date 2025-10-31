package entity

type StateType int

const (
	StateNone StateType = iota
	StateAwaitingVerificationCode
)

type InlineKeyboard struct {
	Name    string
	Command string
}
