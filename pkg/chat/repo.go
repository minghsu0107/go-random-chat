package chat

import (
	"context"
	"errors"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

var (
	messagePubSubTopic = "rc_msg"

	messagesPrefix     = "rc:msgs"
	channelPrefix      = "rc:chan"
	channelUsersPrefix = "rc:chanusers"
	onlineUsersPrefix  = "rc:onlineusers"
	seenMessagesPrefix = "rc:seenmsgs"
)

var (
	ErrChannelNotFound       = errors.New("error channel not found")
	ErrChannelOrUserNotFound = errors.New("error channel or user not found")
)

type UserRepo interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	DeleteAllOnlineUsers(ctx context.Context, channelID uint64) error
}

type MessageRepo interface {
	InsertMessage(ctx context.Context, msg *Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *Message) error
	ListMessages(ctx context.Context, channelID uint64) ([]*Message, error)
}

type ChannelRepo interface {
	CreateChannel(ctx context.Context, channelID uint64) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type UserRepoImpl struct {
	r infra.RedisCache
}

func NewUserRepo(r infra.RedisCache) UserRepo {
	return &UserRepoImpl{r}
}
func (repo *UserRepoImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	key := constructKey(channelUsersPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 0)
}
func (repo *UserRepoImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	key := constructKey(channelUsersPrefix, channelID)
	var dummy int
	return repo.r.HGet(ctx, key, strconv.FormatUint(userID, 10), &dummy)
}
func (repo *UserRepoImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	var dummy int
	exist, err := repo.r.Get(ctx, constructKey(channelPrefix, channelID), &dummy)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrChannelNotFound
	}
	key := constructKey(channelUsersPrefix, channelID)
	userMap, err := repo.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for userIDStr := range userMap {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}
func (repo *UserRepoImpl) AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 0)
}
func (repo *UserRepoImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	userKey := strconv.FormatUint(userID, 10)
	return repo.r.HDel(ctx, key, userKey)
}
func (repo *UserRepoImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	var dummy int
	exist, err := repo.r.Get(ctx, constructKey(channelPrefix, channelID), &dummy)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrChannelNotFound
	}
	key := constructKey(onlineUsersPrefix, channelID)
	userMap, err := repo.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for userIDStr := range userMap {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}
func (repo *UserRepoImpl) DeleteAllOnlineUsers(ctx context.Context, channelID uint64) error {
	return repo.r.Delete(ctx, constructKey(onlineUsersPrefix, channelID))
}

type MessageRepoImpl struct {
	r           infra.RedisCache
	p           message.Publisher
	maxMessages int64
}

func NewMessageRepo(config *config.Config, r infra.RedisCache, p message.Publisher) MessageRepo {
	return &MessageRepoImpl{r, p, config.Chat.Message.MaxNum}
}

func (repo *MessageRepoImpl) InsertMessage(ctx context.Context, msg *Message) error {
	cmds := []infra.RedisCmd{
		{
			OpType: infra.RPUSH,
			Payload: infra.RedisRpushPayload{
				Key: constructKey(messagesPrefix, msg.ChannelID),
				Val: msg.Encode(),
			},
		},
		{
			OpType: infra.HSETONE,
			Payload: infra.RedisHsetOnePayload{
				Key:   constructKey(seenMessagesPrefix, msg.ChannelID),
				Field: strconv.FormatUint(msg.MessageID, 10),
				Val:   0,
			},
		},
	}
	return repo.r.ExecPipeLine(ctx, &cmds)
}
func (repo *MessageRepoImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	key := constructKey(seenMessagesPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(messageID, 10), 1)
}
func (repo *MessageRepoImpl) PublishMessage(ctx context.Context, msg *Message) error {
	return repo.p.Publish(messagePubSubTopic, message.NewMessage(
		watermill.NewUUID(),
		msg.Encode(),
	))
}
func (repo *MessageRepoImpl) ListMessages(ctx context.Context, channelID uint64) ([]*Message, error) {
	var dummy int
	exist, err := repo.r.Get(ctx, constructKey(channelPrefix, channelID), &dummy)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrChannelNotFound
	}

	messageStrs, err := repo.r.LRange(ctx, constructKey(messagesPrefix, channelID), -repo.maxMessages, -1)
	if err != nil {
		return nil, err
	}

	var messages []*Message
	if len(messageStrs) == 0 {
		return messages, nil
	}

	var messageIDStrs []string
	for _, messageStr := range messageStrs {
		message, _ := DecodeToMessage([]byte(messageStr))
		messages = append(messages, message)
		messageIDStrs = append(messageIDStrs, strconv.FormatUint(message.MessageID, 10))
	}

	seenStatuses, err := repo.r.HMGet(ctx, constructKey(seenMessagesPrefix, channelID), messageIDStrs)
	if err != nil {
		return nil, err
	}
	for idx, seenStatus := range seenStatuses {
		switch s := seenStatus.(type) {
		case string:
			val, _ := strconv.ParseInt(s, 10, 64)
			messages[idx].Seen = (val == 1)
		case nil:
			messages[idx].Seen = false
		}
	}
	return messages, nil
}

type ChannelRepoImpl struct {
	r infra.RedisCache
}

func NewChannelRepo(r infra.RedisCache) ChannelRepo {
	return &ChannelRepoImpl{r}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context, channelID uint64) (*Channel, error) {
	if err := repo.r.Set(ctx, constructKey(channelPrefix, channelID), 0); err != nil {
		return nil, err
	}
	return &Channel{
		ID: channelID,
	}, nil
}
func (repo *ChannelRepoImpl) IsChannelExist(ctx context.Context, channelID uint64) (bool, error) {
	var dummy int
	return repo.r.Get(ctx, constructKey(channelPrefix, channelID), &dummy)
}
func (repo *ChannelRepoImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	cmds := []infra.RedisCmd{
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: constructKey(channelPrefix, channelID),
			},
		},
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: constructKey(channelUsersPrefix, channelID),
			},
		},
	}
	return repo.r.ExecPipeLine(ctx, &cmds)
}

func constructKey(prefix string, id uint64) string {
	return prefix + ":" + strconv.FormatUint(id, 10)
}
