package chat

import "errors"

var (
	ErrUserNotFound           = errors.New("error user not found")
	ErrChannelOrUserNotFound  = errors.New("error channel or user not found")
	ErrExceedMessageNumLimits = errors.New("error exceed max number of messages")
)
