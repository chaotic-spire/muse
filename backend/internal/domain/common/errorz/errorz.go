package errorz

import "errors"

var (
	UserAlreadyExists   = errors.New("user already exists")
	AuthHeaderIsEmpty   = errors.New("auth header is empty")
	TelegramUserEmpty   = errors.New("telegram user empty")
	Forbidden           = errors.New("forbidden")
	InvalidHash         = errors.New("invalid hash")
	IncompatibleVersion = errors.New("incompatible version")
)
