package randomchat

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/sony/sonyflake"
)

// IDGenerator is the inteface for generatring unique ID
type IDGenerator interface {
	NextID() (uint64, error)
}

func NewSonyFlake() (IDGenerator, error) {
	var st sonyflake.Settings
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return nil, errors.New("sonyflake not created")
	}
	return sf, nil
}

func DecodeToMatchResult(data []byte) (*MatchResult, error) {
	var result MatchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func DecodeToMessagePresenter(data []byte) (*MessagePresenter, error) {
	var msg MessagePresenter
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func DecodeToMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func getServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
