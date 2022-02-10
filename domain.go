package main

import (
	"encoding/json"
	"strconv"
)

const (
	EventText = iota
	EventAction
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
	Event     int    `json:"event"`
	ChannelID uint64 `json:"channel_id"`
	UserID    uint64 `json:"user_id"`
	Payload   string `json:"payload"`
	Time      string `json:"time"`
}

type Channel struct {
	ID uint64
}

type User struct {
	ID   uint64
	Name string
}

type MatchResult struct {
	Matched   bool
	UserID    uint64
	PeerID    uint64
	ChannelID uint64
}

func (r *MatchResult) Encode() []byte {
	result, _ := json.Marshal(r)
	return result
}
func (r *MatchResult) ToPresenter() *MatchResultPresenter {
	return &MatchResultPresenter{
		ChannelID: strconv.FormatUint(r.ChannelID, 10),
	}
}

func (m *Message) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}

func (m *Message) ToPresenter() *MessagePresenter {
	return &MessagePresenter{
		Event:     m.Event,
		ChannelID: strconv.FormatUint(m.ChannelID, 10),
		UserID:    strconv.FormatUint(m.UserID, 10),
		Payload:   m.Payload,
		Time:      m.Time,
	}
}
