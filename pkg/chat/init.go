package chat

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func InitializeRouter() (*Router, error) {
	svr := gin.Default()
	svr.Use(MaxAllowed(maxAllowedConns))

	MelodyMatch = melody.New()
	MelodyChat = melody.New()
	maxMessageSize, err := strconv.ParseInt(getenv("MAX_MSG_SIZE_BYTE", "4096"), 10, 64)
	if err != nil {
		panic(err)
	}
	MelodyChat.Config.MaxMessageSize = maxMessageSize

	redisClient, err := NewRedisClient()
	if err != nil {
		return nil, err
	}

	redisCache := NewRedisCache(redisClient)

	kafkaPub, err := NewKafkaPublisher()
	if err != nil {
		return nil, err
	}
	kafkaSub, err := NewKafkaSubscriber()
	if err != nil {
		return nil, err
	}

	userRepo := NewRedisUserRepo(redisCache)
	msgRepo := NewMessageRepo(redisCache, kafkaPub)
	chanRepo := NewRedisChannelRepo(redisCache)
	matchRepo := NewMatchingRepo(redisCache, kafkaPub)

	matchSubscriber, err := NewMatchSubscriber(MelodyMatch, userRepo, kafkaSub)
	if err != nil {
		return nil, err
	}
	msgSubscriber, err := NewMessageSubscriber(kafkaSub, MelodyChat)
	if err != nil {
		return nil, err
	}

	sf, err := NewSonyFlake()
	if err != nil {
		return nil, err
	}

	userSvc := NewUserService(userRepo, sf)
	msgSvc := NewMessageService(msgRepo, userRepo, sf)
	matchSvc := NewMatchingService(matchRepo, chanRepo, sf)
	chanSvc := NewChannelService(chanRepo, userRepo)

	return NewRouter(svr, MelodyMatch, MelodyChat, matchSubscriber, msgSubscriber, userSvc, msgSvc, matchSvc, chanSvc), nil
}
