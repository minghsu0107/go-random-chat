package chat

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"gopkg.in/olahol/melody.v1"
)

func InitializeRouter() (*Router, error) {
	svr := gin.Default()
	svr.Use(common.MaxAllowed(maxAllowedConns))

	MelodyMatch = melody.New()
	MelodyChat = melody.New()
	maxMessageSize, err := strconv.ParseInt(common.Getenv("MAX_MSG_SIZE_BYTE", "4096"), 10, 64)
	if err != nil {
		panic(err)
	}
	MelodyChat.Config.MaxMessageSize = maxMessageSize

	redisClient, err := NewRedisClient()
	if err != nil {
		return nil, err
	}

	msgSubscriber := NewMessageSubscriber(redisClient, MelodyChat)

	redisCache := NewRedisCache(redisClient)

	userRepo := NewRedisUserRepo(redisCache)
	msgRepo := NewRedisMessageRepo(redisCache)
	chanRepo := NewRedisChannelRepo(redisCache)
	matchRepo := NewRedisMatchingRepo(redisCache)

	matchSubscriber := NewMatchSubscriber(redisClient, MelodyMatch, userRepo)

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
