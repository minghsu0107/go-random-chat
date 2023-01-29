package chat

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"gopkg.in/olahol/melody.v1"
)

type MessageSubscriber struct {
	subscriberID string
	router       *message.Router
	sub          message.Subscriber
	m            MelodyChatConn
}

func NewMessageSubscriber(name string, config *config.Config, sub message.Subscriber, m MelodyChatConn) (*MessageSubscriber, error) {
	router, err := infra.NewBrokerRouter(name)
	if err != nil {
		return nil, err
	}
	return &MessageSubscriber{
		subscriberID: config.Chat.Subscriber.Id,
		router:       router,
		sub:          sub,
		m:            m,
	}, nil
}

func (s *MessageSubscriber) HandleMessage(msg *message.Message) error {
	message, err := DecodeToMessage([]byte(msg.Payload))
	if err != nil {
		return err
	}
	return s.sendMessage(context.Background(), message)
}

func (s *MessageSubscriber) RegisterHandler() {
	s.router.AddNoPublisherHandler(
		"randomchat_message_handler",
		s.subscriberID,
		s.sub,
		s.HandleMessage,
	)
}

func (s *MessageSubscriber) Run() error {
	return s.router.Run(context.Background())
}

func (s *MessageSubscriber) GracefulStop() error {
	return s.router.Close()
}

func (s *MessageSubscriber) sendMessage(ctx context.Context, message *Message) error {
	return s.m.BroadcastFilter(message.ToPresenter().Encode(), func(sess *melody.Session) bool {
		channelID, exist := sess.Get(sessCidKey)
		if !exist {
			return false
		}
		return message.ChannelID == (channelID.(uint64))
	})
}
