package chat

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

var maxMessages int64
var (
	matchPubSubTopic   = "rc_match"
	messagePubSubTopic = "rc_msg"

	messagesPrefix     = "rc:msgs"
	channelPrefix      = "rc:chan"
	userWaitList       = "rc:userwait"
	userPrefix         = "rc:user"
	channelUsersPrefix = "rc:chanusers"
	onlineUsersPrefix  = "rc:onlineusers"
	seenMessagesPrefix = "rc:seenmsgs"
)

var (
	ErrUserNotFound          = errors.New("error user not found")
	ErrChannelNotFound       = errors.New("error channel not found")
	ErrChannelOrUserNotFound = errors.New("error channel or user not found")
)

func init() {
	var err error
	maxMessages, err = strconv.ParseInt(getenv("MAX_MSGS", "500"), 10, 64)
	if err != nil {
		panic(err)
	}
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
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

type MatchingRepo interface {
	PopOrPushWaitList(ctx context.Context, userID uint64) (bool, uint64, error)
	PublishMatchResult(ctx context.Context, result *MatchResult) error
	RemoveFromWaitList(ctx context.Context, userID uint64) error
}

type RedisUserRepo struct {
	r RedisCache
}

func NewRedisUserRepo(r RedisCache) UserRepo {
	return &RedisUserRepo{r}
}
func (repo *RedisUserRepo) CreateUser(ctx context.Context, user *User) (*User, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = repo.r.Set(ctx, constructKey(userPrefix, user.ID), data)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
func (repo *RedisUserRepo) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	key := constructKey(userPrefix, userID)
	var user User
	exist, err := repo.r.Get(ctx, key, &user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
func (repo *RedisUserRepo) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	key := constructKey(channelUsersPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 0)
}
func (repo *RedisUserRepo) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	key := constructKey(channelUsersPrefix, channelID)
	var dummy int
	return repo.r.HGet(ctx, key, strconv.FormatUint(userID, 10), &dummy)
}
func (repo *RedisUserRepo) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
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
func (repo *RedisUserRepo) AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 0)
}
func (repo *RedisUserRepo) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	userKey := strconv.FormatUint(userID, 10)
	return repo.r.HDel(ctx, key, userKey)
}
func (repo *RedisUserRepo) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
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
func (repo *RedisUserRepo) DeleteAllOnlineUsers(ctx context.Context, channelID uint64) error {
	return repo.r.Delete(ctx, constructKey(onlineUsersPrefix, channelID))
}

type MessageRepoImpl struct {
	r RedisCache
	p message.Publisher
}

func NewMessageRepo(r RedisCache, p message.Publisher) MessageRepo {
	return &MessageRepoImpl{r, p}
}

func (repo *MessageRepoImpl) InsertMessage(ctx context.Context, msg *Message) error {
	cmds := []RedisCmd{
		{
			OpType: RPUSH,
			Payload: RedisRpushPayload{
				Key: constructKey(messagesPrefix, msg.ChannelID),
				Val: msg.Encode(),
			},
		},
		{
			OpType: HSETONE,
			Payload: RedisHsetOnePayload{
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

	messageStrs, err := repo.r.LRange(ctx, constructKey(messagesPrefix, channelID), -maxMessages, -1)
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

type RedisChannelRepo struct {
	r RedisCache
}

func NewRedisChannelRepo(r RedisCache) ChannelRepo {
	return &RedisChannelRepo{r}
}

func (repo *RedisChannelRepo) CreateChannel(ctx context.Context, channelID uint64) (*Channel, error) {
	if err := repo.r.Set(ctx, constructKey(channelPrefix, channelID), 0); err != nil {
		return nil, err
	}
	return &Channel{
		ID: channelID,
	}, nil
}
func (repo *RedisChannelRepo) IsChannelExist(ctx context.Context, channelID uint64) (bool, error) {
	var dummy int
	return repo.r.Get(ctx, constructKey(channelPrefix, channelID), &dummy)
}
func (repo *RedisChannelRepo) DeleteChannel(ctx context.Context, channelID uint64) error {
	cmds := []RedisCmd{
		{
			OpType: DELETE,
			Payload: RedisDeletePayload{
				Key: constructKey(channelPrefix, channelID),
			},
		},
		{
			OpType: DELETE,
			Payload: RedisDeletePayload{
				Key: constructKey(channelUsersPrefix, channelID),
			},
		},
	}
	return repo.r.ExecPipeLine(ctx, &cmds)
}

type MatchingRepoImpl struct {
	r RedisCache
	p message.Publisher
}

func NewMatchingRepo(r RedisCache, p message.Publisher) MatchingRepo {
	return &MatchingRepoImpl{r, p}
}
func (repo *MatchingRepoImpl) PopOrPushWaitList(ctx context.Context, userID uint64) (bool, uint64, error) {
	match, peerIDStr, err := repo.r.ZPopMinOrAddOne(ctx, userWaitList, float64(time.Now().Unix()), userID)
	if err != nil {
		return false, 0, err
	}
	if !match {
		return false, 0, nil
	}
	peerID, err := strconv.ParseUint(peerIDStr, 10, 64)
	if err != nil {
		return false, 0, err
	}
	return true, peerID, nil
}
func (repo *MatchingRepoImpl) RemoveFromWaitList(ctx context.Context, userID uint64) error {
	return repo.r.ZRemOne(ctx, userWaitList, userID)
}
func (repo *MatchingRepoImpl) PublishMatchResult(ctx context.Context, result *MatchResult) error {
	return repo.p.Publish(matchPubSubTopic, message.NewMessage(
		watermill.NewUUID(),
		result.Encode(),
	))
}

func constructKey(prefix string, id uint64) string {
	return prefix + ":" + strconv.FormatUint(id, 10)
}
