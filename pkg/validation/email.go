package validation

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyEmail   = errors.New("email cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmptyEmail
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}
