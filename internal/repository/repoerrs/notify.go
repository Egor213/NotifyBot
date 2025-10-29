package repoerrs

import "errors"

var (
	ErrCreateSetting = errors.New("failed to create notification setting")
	ErrDeleteSetting = errors.New("failed to delete notification setting")
	ErrGetSettings   = errors.New("failed to get notification settings")
)
