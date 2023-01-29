package forwarder

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

var (
	forwardPrefix = "rc:forward"
)

type Subscribers map[string]struct{}

type ForwardRepo interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
	GetSubscribers(ctx context.Context, channelID uint64) (Subscribers, error)
	ForwardMessage(ctx context.Context, msg *chat.Message, subscribers Subscribers) error
}

type ForwardRepoImpl struct {
	r infra.RedisCache
	p message.Publisher
}

func NewForwardRepo(r infra.RedisCache, p message.Publisher) ForwardRepo {
	return &ForwardRepoImpl{r, p}
}

func (repo *ForwardRepoImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	key := constructKey(forwardPrefix, channelID)
	return repo.r.HSet(ctx, key, strconv.FormatUint(userID, 10), subscriber)
}

func (repo *ForwardRepoImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	key := constructKey(forwardPrefix, channelID)
	return repo.r.HDel(ctx, key, strconv.FormatUint(userID, 10))
}

func (repo *ForwardRepoImpl) GetSubscribers(ctx context.Context, channelID uint64) (Subscribers, error) {
	key := constructKey(forwardPrefix, channelID)
	sessionMap, err := repo.r.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	subscribers := make(Subscribers)
	for _, subscriber := range sessionMap {
		subscribers[subscriber] = struct{}{}
	}
	return subscribers, nil
}

func (repo *ForwardRepoImpl) ForwardMessage(ctx context.Context, msg *chat.Message, subscribers Subscribers) error {
	var err error
	for subscriber := range subscribers {
		err = repo.p.Publish(subscriber, message.NewMessage(
			watermill.NewUUID(),
			msg.Encode(),
		))
		if err != nil {
			return err
		}
	}
	return nil
}

func constructKey(prefix string, id uint64) string {
	return common.Join(prefix, ":", strconv.FormatUint(id, 10))
}
