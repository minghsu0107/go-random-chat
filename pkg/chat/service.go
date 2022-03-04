package chat

import (
	"context"
	"strconv"
	"time"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type MessageService interface {
	BroadcastTextMessage(ctx context.Context, channelID, userID uint64, payload string) error
	BroadcastConnectMessage(ctx context.Context, channelID, userID uint64) error
	BroadcastActionMessage(ctx context.Context, channelID, userID uint64, action Action) error
	BroadcastFileMessage(ctx context.Context, channelID, userID uint64, payload string) error
	MarkMessageSeen(ctx context.Context, channelID, userID, messageID uint64) error
	ListMessages(ctx context.Context, channelID uint64) ([]*Message, error)
}

type MatchingService interface {
	Match(ctx context.Context, userID uint64) (*MatchResult, error)
	BroadcastMatchResult(ctx context.Context, result *MatchResult) error
	RemoveUserFromWaitList(ctx context.Context, userID uint64) error
}

type UserService interface {
	CreateUser(ctx context.Context, userName string) (*User, error)
	GetUser(ctx context.Context, uid uint64) (*User, error)
	AddUserToChannel(ctx context.Context, channelID, userID uint64) error
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type ChannelService interface {
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type MessageServiceImpl struct {
	msgRepo  MessageRepo
	userRepo UserRepo
	sf       IDGenerator
}

func NewMessageService(msgRepo MessageRepo, userRepo UserRepo, sf IDGenerator) MessageService {
	return &MessageServiceImpl{msgRepo, userRepo, sf}
}
func (svc *MessageServiceImpl) BroadcastTextMessage(ctx context.Context, channelID, userID uint64, payload string) error {
	messageID, err := svc.sf.NextID()
	if err != nil {
		return err
	}
	msg := Message{
		MessageID: messageID,
		Event:     EventText,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   payload,
		Time:      time.Now().UnixMilli(),
	}
	if err := svc.msgRepo.InsertMessage(ctx, &msg); err != nil {
		return err
	}
	return svc.msgRepo.PublishMessage(ctx, &msg)
}
func (svc *MessageServiceImpl) BroadcastConnectMessage(ctx context.Context, channelID, userID uint64) error {
	onnlineUserIDs, err := svc.userRepo.GetOnlineUserIDs(context.Background(), channelID)
	if err != nil {
		return err
	}
	if len(onnlineUserIDs) == 1 {
		return svc.BroadcastActionMessage(ctx, channelID, userID, WaitingMessage)
	}
	return svc.BroadcastActionMessage(ctx, channelID, userID, JoinedMessage)
}
func (svc *MessageServiceImpl) BroadcastActionMessage(ctx context.Context, channelID, userID uint64, action Action) error {
	eventMessageID, err := svc.sf.NextID()
	if err != nil {
		return err
	}
	msg := Message{
		MessageID: eventMessageID,
		Event:     EventAction,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   string(action),
		Time:      time.Now().UnixMilli(),
	}
	return svc.msgRepo.PublishMessage(ctx, &msg)
}
func (svc *MessageServiceImpl) BroadcastFileMessage(ctx context.Context, channelID, userID uint64, payload string) error {
	messageID, err := svc.sf.NextID()
	if err != nil {
		return err
	}
	msg := Message{
		MessageID: messageID,
		Event:     EventFile,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   payload,
		Time:      time.Now().UnixMilli(),
	}
	if err := svc.msgRepo.InsertMessage(ctx, &msg); err != nil {
		return err
	}
	return svc.msgRepo.PublishMessage(ctx, &msg)
}
func (svc *MessageServiceImpl) MarkMessageSeen(ctx context.Context, channelID, userID, messageID uint64) error {
	if err := svc.msgRepo.MarkMessageSeen(ctx, channelID, messageID); err != nil {
		return err
	}
	eventMessageID, err := svc.sf.NextID()
	if err != nil {
		return err
	}
	msg := Message{
		MessageID: eventMessageID,
		Event:     EventSeen,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   strconv.FormatUint(messageID, 10),
		Seen:      true,
		Time:      time.Now().UnixMilli(),
	}
	return svc.msgRepo.PublishMessage(ctx, &msg)
}
func (svc *MessageServiceImpl) ListMessages(ctx context.Context, channelID uint64) ([]*Message, error) {
	return svc.msgRepo.ListMessages(ctx, channelID)
}

type MatchingServiceImpl struct {
	matchRepo MatchingRepo
	chanRepo  ChannelRepo
	sf        IDGenerator
}

func NewMatchingService(matchRepo MatchingRepo, chanRepo ChannelRepo, sf IDGenerator) MatchingService {
	return &MatchingServiceImpl{matchRepo, chanRepo, sf}
}
func (svc *MatchingServiceImpl) Match(ctx context.Context, userID uint64) (*MatchResult, error) {
	matched, peerID, err := svc.matchRepo.PopOrPushWaitList(ctx, userID)
	if err != nil {
		return nil, err
	}
	if matched {
		newChannelID, err := svc.sf.NextID()
		if err != nil {
			return nil, err
		}
		_, err = svc.chanRepo.CreateChannel(ctx, newChannelID)
		if err != nil {
			return nil, err
		}
		accessToken, err := common.NewJWT(newChannelID)
		if err != nil {
			return nil, err
		}
		return &MatchResult{
			Matched:     true,
			UserID:      userID,
			PeerID:      peerID,
			ChannelID:   newChannelID,
			AccessToken: accessToken,
		}, nil
	}
	return &MatchResult{
		Matched: false,
	}, nil
}
func (svc *MatchingServiceImpl) BroadcastMatchResult(ctx context.Context, result *MatchResult) error {
	return svc.matchRepo.PublishMatchResult(ctx, result)
}
func (svc *MatchingServiceImpl) RemoveUserFromWaitList(ctx context.Context, userID uint64) error {
	return svc.matchRepo.RemoveFromWaitList(ctx, userID)
}

type UserServiceImpl struct {
	userRepo UserRepo
	sf       IDGenerator
}

func NewUserService(userRepo UserRepo, sf IDGenerator) UserService {
	return &UserServiceImpl{userRepo, sf}
}
func (svc *UserServiceImpl) CreateUser(ctx context.Context, userName string) (*User, error) {
	userID, err := svc.sf.NextID()
	if err != nil {
		return nil, err
	}
	return svc.userRepo.CreateUser(ctx, &User{
		ID:   userID,
		Name: userName,
	})
}
func (svc *UserServiceImpl) GetUser(ctx context.Context, uid uint64) (*User, error) {
	return svc.userRepo.GetUserByID(ctx, uid)
}
func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {
	return svc.userRepo.AddUserToChannel(ctx, channelID, userID)
}
func (svc *UserServiceImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	return svc.userRepo.IsChannelUserExist(ctx, channelID, userID)
}
func (svc *UserServiceImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	return svc.userRepo.GetChannelUserIDs(ctx, channelID)
}
func (svc *UserServiceImpl) AddOnlineUser(ctx context.Context, channelID, userID uint64) error {
	return svc.userRepo.AddOnlineUser(ctx, channelID, userID)
}
func (svc *UserServiceImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	return svc.userRepo.DeleteOnlineUser(ctx, channelID, userID)
}
func (svc *UserServiceImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	return svc.userRepo.GetOnlineUserIDs(ctx, channelID)
}

type ChannelServiceImpl struct {
	chanRepo ChannelRepo
	userRepo UserRepo
}

func NewChannelService(chanRepo ChannelRepo, userRepo UserRepo) ChannelService {
	return &ChannelServiceImpl{chanRepo, userRepo}
}
func (svc *ChannelServiceImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := svc.chanRepo.DeleteChannel(ctx, channelID); err != nil {
		return err
	}
	if err := svc.userRepo.DeleteAllOnlineUsers(ctx, channelID); err != nil {
		return err
	}
	return nil
}