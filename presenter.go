package randomchat

import (
	"encoding/json"
	"errors"
	"strconv"
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

type MessagePresenter struct {
	MessageID string `json:"message_id"`
	Event     int    `json:"event"`
	UserID    string `json:"user_id"`
	Payload   string `json:"payload"`
	Seen      bool   `json:"seen"`
	Time      int64  `json:"time"`
}

type ChannelPresenter struct {
	ID string `json:"id"`
}

type ChannelCreation struct {
	UserID string `json:"user_id"`
}

type UserPresenter struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

type UserIDsPresenter struct {
	UserIDs []string `json:"user_ids"`
}

type MessagesPresenter struct {
	Messages []MessagePresenter `json:"messages"`
}

type MatchResultPresenter struct {
	AccessToken string `json:"access_token"`
}

func (m *MatchResultPresenter) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}

func (m *MessagePresenter) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}

func (m *MessagePresenter) ToMessage(accessToken string) (*Message, error) {
	authResult, err := Auth(&AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	if authResult.Expired {
		return nil, ErrTokenExpired
	}
	channelID := authResult.ChannelID
	userID, err := strconv.ParseUint(m.UserID, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Message{
		Event:     m.Event,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   m.Payload,
		Time:      m.Time,
	}, nil
}
