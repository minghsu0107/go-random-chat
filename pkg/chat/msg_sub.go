package chat

import (
	"context"
	"time"

	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"

	retry "github.com/avast/retry-go"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

type MessageSubscriber interface {
	Subscribe() error
	Close()
}

type MessageSubscriberImpl struct {
	client       redis.UniversalClient
	m            MelodyChatConn
	numberWorker int
	pool         *common.Pool
}

func NewMessageSubscriber(config *config.Config, client redis.UniversalClient, m MelodyChatConn) MessageSubscriber {
	return &MessageSubscriberImpl{
		client:       client,
		m:            m,
		numberWorker: config.Chat.Message.Worker,
	}
}

func (s *MessageSubscriberImpl) Subscribe() error {
	ctx := context.Background()
	pubsub := s.client.Subscribe(ctx, messagePubSubTopic)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}
	channel := pubsub.Channel()
	s.pool = common.NewPool(ctx, common.Option{NumberWorker: s.numberWorker})
	s.pool.Start()

	for msg := range channel {
		message, err := DecodeToMessage([]byte(msg.Payload))
		if err != nil {
			log.Error(err)
			continue
		}
		s.pool.Do(s.sendMessage(ctx, message))
	}
	return nil
}

func (s *MessageSubscriberImpl) Close() {
	s.pool.Stop()
}

func (s *MessageSubscriberImpl) sendMessage(ctx context.Context, message *Message) *common.Task {
	return common.NewTask(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, retry.Do(
			func() error {
				return s.m.BroadcastFilter(message.ToPresenter().Encode(), func(sess *melody.Session) bool {
					channelID, exist := sess.Get(sessCidKey)
					if !exist {
						return false
					}
					return message.ChannelID == (channelID.(uint64))
				})
			},
			retry.Attempts(3),
			retry.DelayType(retry.RandomDelay),
			retry.MaxJitter(10*time.Millisecond),
		)
	})
}
