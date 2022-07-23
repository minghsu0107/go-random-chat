package chat

import (
	"context"
	"errors"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gocql/gocql"
	"github.com/minghsu0107/go-random-chat/pkg/config"
)

var (
	messagePubSubTopic = "rc_msg"
)

var (
	ErrChannelOrUserNotFound = errors.New("error channel or user not found")
)

type UserRepo interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type MessageRepo interface {
	InsertMessage(ctx context.Context, msg *Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *Message) error
	ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*Message, string, error)
}

type ChannelRepo interface {
	CreateChannel(ctx context.Context, channelID uint64) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type UserRepoImpl struct {
	s *gocql.Session
}

func NewUserRepo(s *gocql.Session) UserRepo {
	return &UserRepoImpl{s}
}
func (repo *UserRepoImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	if err := repo.s.Query("INSERT INTO channels (id, user_id) VALUES (?, ?)",
		channelID, userID).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}
func (repo *UserRepoImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	iter := repo.s.Query("SELECT user_id FROM channels WHERE id = ?", channelID).WithContext(ctx).Iter()
	var userIDs []uint64
	var userID uint64
	for iter.Scan(&userID) {
		userIDs = append(userIDs, userID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return userIDs, nil
}

type MessageRepoImpl struct {
	s           *gocql.Session
	p           message.Publisher
	maxMessages int64
	pagination  int
}

func NewMessageRepo(config *config.Config, s *gocql.Session, p message.Publisher) MessageRepo {
	return &MessageRepoImpl{s, p, config.Chat.Message.MaxNum, config.Chat.Message.PaginationNum}
}

func (repo *MessageRepoImpl) InsertMessage(ctx context.Context, msg *Message) error {
	if err := repo.s.Query("INSERT INTO messages (id, event, channel_id, user_id, payload, seen, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
		msg.MessageID,
		msg.Event,
		msg.ChannelID,
		msg.UserID,
		msg.Payload,
		false,
		msg.Time).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}
func (repo *MessageRepoImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	if err := repo.s.Query("UPDATE messages SET seen = ? WHERE channel_id = ? AND id = ?", true, channelID, messageID).
		WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}
func (repo *MessageRepoImpl) PublishMessage(ctx context.Context, msg *Message) error {
	return repo.p.Publish(messagePubSubTopic, message.NewMessage(
		watermill.NewUUID(),
		msg.Encode(),
	))
}
func (repo *MessageRepoImpl) ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*Message, string, error) {
	var messages []*Message
	pageState := []byte(pageStateStr)
	iter := repo.s.Query(`SELECT id, event, channel_id, user_id, payload, seen, timestamp FROM messages WHERE channel_id = ?`, channelID).
		WithContext(ctx).PageSize(repo.pagination).PageState(pageState).Iter()
	nextPageStateStr := string(iter.PageState())
	scanner := iter.Scanner()
	var err error
	for scanner.Next() {
		var message Message
		if err = scanner.Scan(
			&message.MessageID,
			&message.Event,
			&message.ChannelID,
			&message.UserID,
			&message.Payload,
			&message.Seen,
			&message.Time); err != nil {
			return nil, "", err
		}
		messages = append(messages, &message)
	}
	err = scanner.Err()
	if err != nil {
		return nil, "", err
	}
	return messages, nextPageStateStr, nil
}

type ChannelRepoImpl struct {
	s *gocql.Session
}

func NewChannelRepo(s *gocql.Session) ChannelRepo {
	return &ChannelRepoImpl{s}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context, channelID uint64) (*Channel, error) {
	if err := repo.s.Query("INSERT INTO channels (id, user_id) VALUES (?, ?)",
		channelID, 0).WithContext(ctx).Exec(); err != nil {
		return nil, err
	}
	return &Channel{
		ID: channelID,
	}, nil
}
func (repo *ChannelRepoImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := repo.s.Query("DELETE FROM channels WHERE id = ?", channelID).
		WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func constructKey(prefix string, id uint64) string {
	return prefix + ":" + strconv.FormatUint(id, 10)
}
