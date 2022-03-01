package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func InitializeRouter() (*Router, error) {
	svr := gin.Default()
	svr.Use(MaxAllowed(maxAllowedConns))

	MelodyMatch = melody.New()
	MelodyChat = melody.New()

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
