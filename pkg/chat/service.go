package chat

import (
	"context"
	"fmt"
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
	InsertMessage(ctx context.Context, msg *Message) error
	PublishMessage(ctx context.Context, msg *Message) error
	ListMessages(ctx context.Context, channelID uint64, pageState string) ([]*Message, string, error)
}

type UserService interface {
	AddUserToChannel(ctx context.Context, channelID, userID uint64) error
	GetUser(ctx context.Context, userID uint64) (*User, error)
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type ChannelService interface {
	CreateChannel(ctx context.Context) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type ForwardService interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
}

type MessageServiceImpl struct {
	msgRepo  MessageRepoCache
	userRepo UserRepoCache
	sf       common.IDGenerator
}

func NewMessageService(msgRepo MessageRepoCache, userRepo UserRepoCache, sf common.IDGenerator) MessageService {
	return &MessageServiceImpl{msgRepo, userRepo, sf}
}
func (svc *MessageServiceImpl) BroadcastTextMessage(ctx context.Context, channelID, userID uint64, payload string) error {
	messageID, err := svc.sf.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for text message: %w", err)
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
		return fmt.Errorf("error broadcast text message: %w", err)
	}
	if err := svc.PublishMessage(ctx, &msg); err != nil {
		return fmt.Errorf("error broadcast text message: %w", err)
	}
	return nil
}
func (svc *MessageServiceImpl) BroadcastConnectMessage(ctx context.Context, channelID, userID uint64) error {
	onnlineUserIDs, err := svc.userRepo.GetOnlineUserIDs(context.Background(), channelID)
	if err != nil {
		return fmt.Errorf("error get online user ids from channel %d: %w", channelID, err)
	}
	if len(onnlineUserIDs) == 1 {
		return svc.BroadcastActionMessage(ctx, channelID, userID, WaitingMessage)
	}
	return svc.BroadcastActionMessage(ctx, channelID, userID, JoinedMessage)
}
func (svc *MessageServiceImpl) BroadcastActionMessage(ctx context.Context, channelID, userID uint64, action Action) error {
	eventMessageID, err := svc.sf.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for action message: %w", err)
	}
	msg := Message{
		MessageID: eventMessageID,
		Event:     EventAction,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   string(action),
		Time:      time.Now().UnixMilli(),
	}
	if err := svc.PublishMessage(ctx, &msg); err != nil {
		return fmt.Errorf("error broadcast action message: %w", err)
	}
	return nil
}
func (svc *MessageServiceImpl) BroadcastFileMessage(ctx context.Context, channelID, userID uint64, payload string) error {
	messageID, err := svc.sf.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for file message: %w", err)
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
		return fmt.Errorf("error broadcast file message: %w", err)
	}
	if err := svc.PublishMessage(ctx, &msg); err != nil {
		return fmt.Errorf("error broadcast file message: %w", err)
	}
	return nil
}
func (svc *MessageServiceImpl) MarkMessageSeen(ctx context.Context, channelID, userID, messageID uint64) error {
	if err := svc.msgRepo.MarkMessageSeen(ctx, channelID, messageID); err != nil {
		return fmt.Errorf("error mark message %d seen in channel %d: %w", messageID, channelID, err)
	}
	eventMessageID, err := svc.sf.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for seen event message: %w", err)
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
	if err := svc.PublishMessage(ctx, &msg); err != nil {
		return fmt.Errorf("error mark message %d seen in channel %d: %w", messageID, channelID, err)
	}
	return nil
}
func (svc *MessageServiceImpl) InsertMessage(ctx context.Context, msg *Message) error {
	if err := svc.msgRepo.InsertMessage(ctx, msg); err != nil {
		return fmt.Errorf("error insert message: %w", err)
	}
	return nil
}
func (svc *MessageServiceImpl) PublishMessage(ctx context.Context, msg *Message) error {
	if err := svc.msgRepo.PublishMessage(ctx, msg); err != nil {
		return fmt.Errorf("error publish message: %w", err)
	}
	return nil
}
func (svc *MessageServiceImpl) ListMessages(ctx context.Context, channelID uint64, pageState string) ([]*Message, string, error) {
	msgs, nextPageState, err := svc.msgRepo.ListMessages(ctx, channelID, pageState)
	if err != nil {
		return nil, "", fmt.Errorf("error list messages in channel %d with page state %s: %w", channelID, pageState, err)
	}
	return msgs, nextPageState, nil
}

type UserServiceImpl struct {
	userRepo UserRepoCache
}

func NewUserService(userRepo UserRepoCache) UserService {
	return &UserServiceImpl{userRepo}
}
func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.AddUserToChannel(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error add user %d to channel %d: %w", userID, channelID, err)
	}
	return nil
}
func (svc *UserServiceImpl) GetUser(ctx context.Context, userID uint64) (*User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", userID, err)
	}
	return user, nil
}
func (svc *UserServiceImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	exist, err := svc.userRepo.IsChannelUserExist(ctx, channelID, userID)
	if err != nil {
		return false, fmt.Errorf("error check user %d in channel %d: %w", userID, channelID, err)
	}
	return exist, nil
}
func (svc *UserServiceImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	users, err := svc.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error get users in channel %d: %w", channelID, err)
	}
	return users, nil
}
func (svc *UserServiceImpl) AddOnlineUser(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.AddOnlineUser(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error add online user %d to channel %d: %w", userID, channelID, err)
	}
	return nil
}
func (svc *UserServiceImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.DeleteOnlineUser(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error delete online user %d from channel %d: %w", userID, channelID, err)
	}
	return nil
}
func (svc *UserServiceImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	users, err := svc.userRepo.GetOnlineUserIDs(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error get online users in channel %d: %w", channelID, err)
	}
	return users, nil
}

type ChannelServiceImpl struct {
	chanRepo ChannelRepoCache
	userRepo UserRepoCache
	sf       common.IDGenerator
}

func NewChannelService(chanRepo ChannelRepoCache, userRepo UserRepoCache, sf common.IDGenerator) ChannelService {
	return &ChannelServiceImpl{chanRepo, userRepo, sf}
}
func (svc *ChannelServiceImpl) CreateChannel(ctx context.Context) (*Channel, error) {
	channelID, err := svc.sf.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID for new channel: %w", err)
	}
	channel, err := svc.chanRepo.CreateChannel(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error create channel %d: %w", channelID, err)
	}
	return channel, nil
}
func (svc *ChannelServiceImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := svc.chanRepo.DeleteChannel(ctx, channelID); err != nil {
		return fmt.Errorf("error delete channel %d: %w", channelID, err)
	}
	return nil
}

type ForwardServiceImpl struct {
	forwardRepo ForwardRepo
}

func NewForwardService(forwardRepo ForwardRepo) ForwardService {
	return &ForwardServiceImpl{forwardRepo}
}

func (svc *ForwardServiceImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	return svc.forwardRepo.RegisterChannelSession(ctx, channelID, userID, subscriber)
}
func (svc *ForwardServiceImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	return svc.forwardRepo.RemoveChannelSession(ctx, channelID, userID)
}
