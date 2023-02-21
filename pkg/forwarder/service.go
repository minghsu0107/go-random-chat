package forwarder

import (
	"context"

	"github.com/minghsu0107/go-random-chat/pkg/chat"
)

type ForwardService interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
	ForwardMessage(ctx context.Context, msg *chat.Message) error
}

type ForwardServiceImpl struct {
	forwardRepo ForwardRepo
}

func NewForwardServiceImpl(forwardRepo ForwardRepo) *ForwardServiceImpl {
	return &ForwardServiceImpl{forwardRepo}
}

func (svc *ForwardServiceImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	return svc.forwardRepo.RegisterChannelSession(ctx, channelID, userID, subscriber)
}

func (svc *ForwardServiceImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	return svc.forwardRepo.RemoveChannelSession(ctx, channelID, userID)
}

func (svc *ForwardServiceImpl) ForwardMessage(ctx context.Context, msg *chat.Message) error {
	subscribers, err := svc.forwardRepo.GetSubscribers(ctx, msg.ChannelID)
	if err != nil {
		return err
	}
	return svc.forwardRepo.ForwardMessage(ctx, msg, subscribers)
}
