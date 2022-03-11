package chat

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

type MatchSubscriber struct {
	m        MelodyMatchConn
	router   *message.Router
	userRepo UserRepo
	sub      message.Subscriber
}

func NewMatchSubscriber(m MelodyMatchConn, userRepo UserRepo, sub message.Subscriber) (*MatchSubscriber, error) {
	router, err := infra.NewMessageRouter()
	if err != nil {
		return nil, err
	}
	return &MatchSubscriber{
		m:        m,
		router:   router,
		userRepo: userRepo,
		sub:      sub,
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
	s.RegisterHandler()
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
			if err := s.userRepo.AddUserToChannel(ctx, result.ChannelID, userID); err != nil {
				log.Error(err)
				return false
			}
			return true
		}
		return false
	})
}
