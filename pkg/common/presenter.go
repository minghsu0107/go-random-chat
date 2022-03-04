package common

import (
	"errors"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrServer       = errors.New("server error")
	ErrUnauthorized = errors.New("unauthorized")
)

// ErrResponse is the error response type
type ErrResponse struct {
	Message string `json:"msg"`
}

// SuccessMessage is the success response type
type SuccessMessage struct {
	Message string `json:"msg" example:"ok"`
}

// OkMsg is the default success response for 200 status code
var OkMsg SuccessMessage = SuccessMessage{
	Message: "ok",
}
