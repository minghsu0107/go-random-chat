package user

import "errors"

var (
	ErrUserNotFound    = errors.New("error user not found")
	ErrSessionNotFound = errors.New("error session not found")
)
