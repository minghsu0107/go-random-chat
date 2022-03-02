package randomchat

import (
	"context"
	"strconv"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

var msgNumberWorker int64

type MessageSubscriber interface {
	Subscribe() error
	Close()
}

type MessageSubscriberImpl struct {
	client redis.UniversalClient
	m      *melody.Melody
	pool   *Pool
}

func init() {
	var err error
	msgNumberWorker, err = strconv.ParseInt(getenv("MSG_NUMBER_WORKER", "4"), 10, 0)
	if err != nil {
		panic(err)
	}
}

func NewMessageSubscriber(client redis.UniversalClient, m *melody.Melody) MessageSubscriber {
	return &MessageSubscriberImpl{
		client: client,
		m:      m,
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
	s.pool = NewPool(ctx, Option{NumberWorker: int(msgNumberWorker)})
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

func (s *MessageSubscriberImpl) sendMessage(ctx context.Context, message *Message) *Task {
	return NewTask(ctx, func(ctx context.Context) (interface{}, error) {
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
