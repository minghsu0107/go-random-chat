package chat

import (
	"context"
	b64 "encoding/base64"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gocql/gocql"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"

	"github.com/go-kit/kit/endpoint"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	forwarderpb "github.com/minghsu0107/go-random-chat/proto/forwarder"
	userpb "github.com/minghsu0107/go-random-chat/proto/user"
)

var (
	MessagePubTopic = "rc.msg.pub"
)

type UserRepo interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type MessageRepo interface {
	InsertMessage(ctx context.Context, msg *Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *Message) error
	ListMessages(ctx context.Context, channelID uint64, pageStateBase64 string) ([]*Message, string, error)
}

type ChannelRepo interface {
	CreateChannel(ctx context.Context, channelID uint64) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type ForwardRepo interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
}

type UserRepoImpl struct {
	s       *gocql.Session
	getUser endpoint.Endpoint
}

func NewUserRepo(s *gocql.Session, userConn *UserClientConn) UserRepo {
	return &UserRepoImpl{
		s: s,
		getUser: transport.NewGrpcEndpoint(
			userConn.Conn,
			"user",
			"user.UserService",
			"GetUser",
			&userpb.GetUserResponse{},
		),
	}
}
func (repo *UserRepoImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	if err := repo.s.Query("INSERT INTO channels (id, user_id) VALUES (?, ?)",
		channelID, userID).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}
func (repo *UserRepoImpl) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	res, err := repo.getUser(ctx, &userpb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	pbUser := res.(*userpb.GetUserResponse)
	if !pbUser.Exist {
		return nil, ErrUserNotFound
	}
	return &User{
		ID:   pbUser.User.Id,
		Name: pbUser.User.Name,
	}, nil
}
func (repo *UserRepoImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	iter := repo.s.Query("SELECT user_id FROM channels WHERE id = ?", channelID).WithContext(ctx).Idempotent(true).Iter()
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
	var messageNum int64
	err := repo.s.Query("SELECT msgnum FROM chanmsg_counters WHERE channel_id = ? LIMIT 1", msg.ChannelID).
		WithContext(ctx).Idempotent(true).Scan(&messageNum)
	if err != nil {
		if err == gocql.ErrNotFound {
			messageNum = 0
		} else {
			return err
		}
	}
	if messageNum >= repo.maxMessages {
		return ErrExceedMessageNumLimits
	}
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
	return repo.s.Query("UPDATE chanmsg_counters SET msgnum = msgnum + 1 WHERE channel_id = ?", msg.ChannelID).WithContext(ctx).Exec()
}
func (repo *MessageRepoImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	if err := repo.s.Query("UPDATE messages SET seen = ? WHERE channel_id = ? AND id = ?", true, channelID, messageID).
		WithContext(ctx).Idempotent(true).Exec(); err != nil {
		return err
	}
	return nil
}
func (repo *MessageRepoImpl) PublishMessage(ctx context.Context, msg *Message) error {
	return repo.p.Publish(MessagePubTopic, message.NewMessage(
		watermill.NewUUID(),
		msg.Encode(),
	))
}
func (repo *MessageRepoImpl) ListMessages(ctx context.Context, channelID uint64, pageStateBase64 string) ([]*Message, string, error) {
	var messages []*Message
	pageState, err := b64.URLEncoding.DecodeString(pageStateBase64)
	if err != nil {
		return nil, "", err
	}
	iter := repo.s.Query(`SELECT id, event, channel_id, user_id, payload, seen, timestamp FROM messages WHERE channel_id = ?`, channelID).
		WithContext(ctx).Idempotent(true).PageSize(repo.pagination).PageState(pageState).Iter()
	nextPageStateBase64 := b64.URLEncoding.EncodeToString(iter.PageState())
	scanner := iter.Scanner()

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
	return messages, nextPageStateBase64, nil
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
	accessToken, err := common.NewJWT(channelID)
	if err != nil {
		return nil, fmt.Errorf("error create JWT: %w", err)
	}
	return &Channel{
		ID:          channelID,
		AccessToken: accessToken,
	}, nil
}
func (repo *ChannelRepoImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := repo.s.Query("DELETE FROM channels WHERE id = ?", channelID).
		WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

type ForwardRepoImpl struct {
	registerChannelSession endpoint.Endpoint
	removeChannelSession   endpoint.Endpoint
}

func NewForwardRepo(forwarderConn *ForwarderClientConn) ForwardRepo {
	return &ForwardRepoImpl{
		registerChannelSession: transport.NewGrpcEndpoint(
			forwarderConn.Conn,
			"forwarder",
			"forwarder.ForwardService",
			"RegisterChannelSession",
			&forwarderpb.RegisterChannelSessionResponse{},
		),
		removeChannelSession: transport.NewGrpcEndpoint(
			forwarderConn.Conn,
			"forwarder",
			"forwarder.ForwardService",
			"RemoveChannelSession",
			&forwarderpb.RemoveChannelSessionResponse{},
		),
	}
}

func (repo *ForwardRepoImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	_, err := repo.registerChannelSession(ctx, &forwarderpb.RegisterChannelSessionRequest{
		ChannelId:  channelID,
		UserId:     userID,
		Subscriber: subscriber,
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *ForwardRepoImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	_, err := repo.removeChannelSession(ctx, &forwarderpb.RemoveChannelSessionRequest{
		ChannelId: channelID,
		UserId:    userID,
	})
	if err != nil {
		return err
	}
	return nil
}
