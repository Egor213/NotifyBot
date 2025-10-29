package repoerrs

import "errors"

var (
	ErrUserExists   = errors.New("пользователь с таким tg_id уже существует")
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrCheckUser    = errors.New("ошибка проверки существования пользователя")
	ErrCreateUser   = errors.New("ошибка создания пользователя")
	ErrGetUser      = errors.New("ошибка получения пользователя")
)
