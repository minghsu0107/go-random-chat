package chat

import (
	"context"
	"strconv"

	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

var (
	channelUsersPrefix = "rc:chanusers"
	onlineUsersPrefix  = "rc:onlineusers"
)

type UserRepoCache interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type MessageRepoCache interface {
	InsertMessage(ctx context.Context, msg *Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *Message) error
	ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*Message, string, error)
}

type ChannelRepoCache interface {
	CreateChannel(ctx context.Context, channelID uint64) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type UserRepoCacheImpl struct {
	r        infra.RedisCache
	userRepo UserRepo
}

func NewUserRepoCacheImpl(r infra.RedisCache, userRepo UserRepo) *UserRepoCacheImpl {
	return &UserRepoCacheImpl{r, userRepo}
}
func (cache *UserRepoCacheImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	if err := cache.userRepo.AddUserToChannel(ctx, channelID, userID); err != nil {
		return nil
	}
	key := constructKey(channelUsersPrefix, channelID)
	return cache.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 1)
}
func (cache *UserRepoCacheImpl) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	return cache.userRepo.GetUserByID(ctx, userID)
}
func (cache *UserRepoCacheImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	key := constructKey(channelUsersPrefix, channelID)
	var dummy int
	var err error
	channelExists, userExists, err := cache.r.HGetIfKeyExists(ctx, key, strconv.FormatUint(userID, 10), &dummy)
	if err != nil {
		return false, err
	}
	if channelExists {
		if !userExists {
			return false, nil
		}
		return true, nil
	}

	mutex := cache.r.GetMutex(common.Join("mutex:", key))
	if err := mutex.LockContext(ctx); err != nil {
		return false, err
	}
	defer func() {
		_, err = mutex.UnlockContext(ctx)
	}()
	channelExists, userExists, err = cache.r.HGetIfKeyExists(ctx, key, strconv.FormatUint(userID, 10), &dummy)
	if err != nil {
		return false, err
	}
	if channelExists {
		if !userExists {
			return false, nil
		}
		return true, nil
	}

	channelUserIDs, err := cache.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return false, err
	}
	channelUserExist := false
	var args []interface{}
	for _, channelUserID := range channelUserIDs {
		if userID == channelUserID {
			channelUserExist = true
		}
		args = append(args, channelUserID, 1)
	}
	if err := cache.r.HSet(ctx, key, args...); err != nil {
		return false, err
	}
	return channelUserExist, nil
}
func (cache *UserRepoCacheImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	key := constructKey(channelUsersPrefix, channelID)
	userMap, err := cache.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	if len(userMap) > 0 {
		for userIDStr := range userMap {
			userID, err := strconv.ParseUint(userIDStr, 10, 64)
			if err != nil {
				return nil, err
			}
			userIDs = append(userIDs, userID)
		}
		return userIDs, nil
	}

	mutex := cache.r.GetMutex(common.Join("mutex:", key))
	if err := mutex.LockContext(ctx); err != nil {
		return nil, err
	}
	defer func() {
		_, err = mutex.UnlockContext(ctx)
	}()
	userMap, err = cache.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(userMap) > 0 {
		for userIDStr := range userMap {
			userID, err := strconv.ParseUint(userIDStr, 10, 64)
			if err != nil {
				return nil, err
			}
			userIDs = append(userIDs, userID)
		}
		return userIDs, nil
	}

	userIDs, err = cache.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return nil, err
	}
	var args []interface{}
	for _, userID := range userIDs {
		args = append(args, userID, 1)
	}
	if err := cache.r.HSet(ctx, key, args...); err != nil {
		return nil, err
	}
	return userIDs, nil
}
func (cache *UserRepoCacheImpl) AddOnlineUser(ctx context.Context, channelID uint64, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	return cache.r.HSet(ctx, key, strconv.FormatUint(userID, 10), 1)
}
func (cache *UserRepoCacheImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	key := constructKey(onlineUsersPrefix, channelID)
	userKey := strconv.FormatUint(userID, 10)
	return cache.r.HDel(ctx, key, userKey)
}
func (cache *UserRepoCacheImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	key := constructKey(onlineUsersPrefix, channelID)
	userMap, err := cache.r.HGetAll(ctx, key)
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

type MessageRepoCacheImpl struct {
	messageRepo MessageRepo
}

func NewMessageRepoCacheImpl(messageRepo MessageRepo) *MessageRepoCacheImpl {
	return &MessageRepoCacheImpl{messageRepo}
}

func (cache *MessageRepoCacheImpl) InsertMessage(ctx context.Context, msg *Message) error {
	return cache.messageRepo.InsertMessage(ctx, msg)
}
func (cache *MessageRepoCacheImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	return cache.messageRepo.MarkMessageSeen(ctx, channelID, messageID)
}
func (cache *MessageRepoCacheImpl) PublishMessage(ctx context.Context, msg *Message) error {
	return cache.messageRepo.PublishMessage(ctx, msg)
}
func (cache *MessageRepoCacheImpl) ListMessages(ctx context.Context, channelID uint64, pageStateStr string) ([]*Message, string, error) {
	return cache.messageRepo.ListMessages(ctx, channelID, pageStateStr)
}

type ChannelRepoCacheImpl struct {
	r           infra.RedisCache
	channelRepo ChannelRepo
}

func NewChannelRepoCacheImpl(r infra.RedisCache, channelRepo ChannelRepo) *ChannelRepoCacheImpl {
	return &ChannelRepoCacheImpl{r, channelRepo}
}

func (cache *ChannelRepoCacheImpl) CreateChannel(ctx context.Context, channelID uint64) (*Channel, error) {
	return cache.channelRepo.CreateChannel(ctx, channelID)
}
func (cache *ChannelRepoCacheImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := cache.channelRepo.DeleteChannel(ctx, channelID); err != nil {
		return err
	}
	cmds := []infra.RedisCmd{
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: constructKey(onlineUsersPrefix, channelID),
			},
		},
		{
			OpType: infra.DELETE,
			Payload: infra.RedisDeletePayload{
				Key: constructKey(channelUsersPrefix, channelID),
			},
		},
	}
	return cache.r.ExecPipeLine(ctx, &cmds)
}

func constructKey(prefix string, id uint64) string {
	return common.Join(prefix, ":", strconv.FormatUint(id, 10))
}
