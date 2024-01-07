package match

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"gopkg.in/olahol/melody.v1"
)

type MatchSubscriber struct {
	m       MelodyMatchConn
	router  *message.Router
	userSvc UserService
	sub     message.Subscriber
}

func NewMatchSubscriber(name string, router *message.Router, m MelodyMatchConn, userSvc UserService, sub message.Subscriber) (*MatchSubscriber, error) {
	return &MatchSubscriber{
		m:       m,
		router:  router,
		userSvc: userSvc,
		sub:     sub,
	}, nil
}

func (s *MatchSubscriber) HandleMatchResult(msg *message.Message) error {
	result, err := DecodeToMatchResult([]byte(msg.Payload))
	if err != nil {
		return err
	}
	return s.sendMatchResult(context.Background(), result)
}

func (s *MatchSubscriber) RegisterHandler() {
	s.router.AddNoPublisherHandler(
		"randomchat_match_result_handler",
		matchPubSubTopic,
		s.sub,
		s.HandleMatchResult,
	)
}

func (s *MatchSubscriber) Run() error {
	return s.router.Run(context.Background())
}

func (s *MatchSubscriber) GracefulStop() error {
	return s.router.Close()
}

func (s *MatchSubscriber) sendMatchResult(ctx context.Context, result *MatchResult) error {
	return s.m.BroadcastFilter(result.ToPresenter().Encode(), func(sess *melody.Session) bool {
		uid, exist := sess.Get(sessUidKey)
		if !exist {
			return false
		}
		userID := uid.(uint64)
		if (userID == result.PeerID) || (userID == result.UserID) {
			if err := s.userSvc.AddUserToChannel(ctx, result.ChannelID, userID); err != nil {
				slog.Error(err.Error())
				return false
			}
			return true
		}
		return false
	})
}
