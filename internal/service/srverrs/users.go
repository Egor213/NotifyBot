package srverrs

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this TG ID already exists")
	ErrUserCreateFailed  = errors.New("failed to create user")
	ErrUserCheckFailed   = errors.New("failed to check user existence")
)
