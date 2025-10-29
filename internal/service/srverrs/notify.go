package srverrs

import "errors"

var (
	ErrSetNotifySettings     = errors.New("failed to set notification setting")
	ErrRemoveNotifySettings  = errors.New("failed to remove notification setting")
	ErrGetNotifySettings     = errors.New("failed to get notification settings")
	ErrNotifySettingNotFound = errors.New("failed notification setting not found")
)
