package match

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

type MatchSubscriber interface {
	Subscribe() error
	Close()
}

type MatchSubscriberImpl struct {
	client       redis.UniversalClient
	m            MelodyMatchConn
	numberWorker int
	pool         *common.Pool
	userSvc      UserService
}

func NewMatchSubscriber(config *config.Config, client redis.UniversalClient, m MelodyMatchConn, userSvc UserService) MatchSubscriber {
	return &MatchSubscriberImpl{
		client:       client,
		m:            m,
		numberWorker: config.Match.Worker,
		userSvc:      userSvc,
	}
}

func (s *MatchSubscriberImpl) Subscribe() error {
	ctx := context.Background()
	pubsub := s.client.Subscribe(ctx, matchPubSubTopic)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}
	channel := pubsub.Channel()
	s.pool = common.NewPool(ctx, common.Option{NumberWorker: s.numberWorker})
	s.pool.Start()

	for msg := range channel {
		result, err := DecodeToMatchResult([]byte(msg.Payload))
		if err != nil {
			log.Error(err)
			continue
		}
		s.pool.Do(s.sendMatchResult(ctx, result))
	}
	return nil
}

func (s *MatchSubscriberImpl) Close() {
	s.pool.Stop()
}

func (s *MatchSubscriberImpl) sendMatchResult(ctx context.Context, result *MatchResult) *common.Task {
	return common.NewTask(ctx, func(ctx context.Context) (interface{}, error) {
		return nil, retry.Do(
			func() error {
				return s.m.BroadcastFilter(result.ToPresenter().Encode(), func(sess *melody.Session) bool {
					uid, exist := sess.Get(sessUidKey)
					if !exist {
						return false
					}
					userID := uid.(uint64)
					if (userID == result.PeerID) || (userID == result.UserID) {
						if err := s.userSvc.AddUserToChannel(ctx, result.ChannelID, userID); err != nil {
							log.Error(err)
							return false
						}
						return true
					}
					return false
				})
			},
			retry.Attempts(3),
			retry.DelayType(retry.RandomDelay),
			retry.MaxJitter(10*time.Millisecond),
		)
	})
}
