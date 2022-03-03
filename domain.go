package randomchat

import (
	"encoding/json"
	"strconv"
)

const (
	EventText = iota
	EventAction
	EventSeen
	EventFile
)

type Action string

var (
	WaitingMessage   Action = "waiting"
	JoinedMessage    Action = "joined"
	IsTypingMessage  Action = "istyping"
	EndTypingMessage Action = "endtyping"
	OfflineMessage   Action = "offline"
	LeavedMessage    Action = "leaved"
)

type Message struct {
	MessageID uint64 `json:"message_id"`
	Event     int    `json:"event"`
	ChannelID uint64 `json:"channel_id"`
	UserID    uint64 `json:"user_id"`
	Payload   string `json:"payload"`
	Seen      bool   `json:"seen"`
	Time      int64  `json:"time"`
}

type Channel struct {
	ID uint64
}

type User struct {
	ID   uint64
	Name string
}

type MatchResult struct {
	Matched     bool
	UserID      uint64
	PeerID      uint64
	ChannelID   uint64
	AccessToken string
}

func (r *MatchResult) Encode() []byte {
	result, _ := json.Marshal(r)
	return result
}
func (r *MatchResult) ToPresenter() *MatchResultPresenter {
	return &MatchResultPresenter{
		AccessToken: r.AccessToken,
	}
}

func (m *Message) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}

func (m *Message) ToPresenter() *MessagePresenter {
	return &MessagePresenter{
		MessageID: strconv.FormatUint(m.MessageID, 10),
		Event:     m.Event,
		UserID:    strconv.FormatUint(m.UserID, 10),
		Payload:   m.Payload,
		Seen:      m.Seen,
		Time:      m.Time,
	}
}
