package forwarder

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
)

type MessageSubscriber struct {
	router     *message.Router
	sub        message.Subscriber
	forwardSvc ForwardService
}

func NewMessageSubscriber(name string, router *message.Router, sub message.Subscriber, forwardSvc ForwardService) (*MessageSubscriber, error) {
	return &MessageSubscriber{
		router:     router,
		sub:        sub,
		forwardSvc: forwardSvc,
	}, nil
}

func (s *MessageSubscriber) HandleMessage(msg *message.Message) error {
	message, err := chat.DecodeToMessage([]byte(msg.Payload))
	if err != nil {
		return err
	}
	return s.forwardSvc.ForwardMessage(msg.Context(), message)
}

func (s *MessageSubscriber) RegisterHandler() {
	s.router.AddNoPublisherHandler(
		"randomchat_message_forwarder",
		chat.MessagePubTopic,
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
