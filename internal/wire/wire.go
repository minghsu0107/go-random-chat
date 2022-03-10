//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/uploader"
	"github.com/minghsu0107/go-random-chat/pkg/web"
)

func InitializeWebRouter() (*web.Router, error) {
	wire.Build(
		config.NewConfig,
		web.NewGinServer,
		web.NewRouter,
	)
	return &web.Router{}, nil
}

func InitializeChatRouter() (*chat.Router, error) {
	wire.Build(
		config.NewConfig,
		chat.NewGinServer,

		chat.NewRedisClient,
		chat.NewRedisCache,

		chat.NewRedisUserRepo,
		chat.NewRedisMessageRepo,
		chat.NewRedisChannelRepo,
		chat.NewRedisMatchingRepo,

		chat.NewMessageSubscriber,
		chat.NewMatchSubscriber,

		chat.NewSonyFlake,

		chat.NewUserService,
		chat.NewMessageService,
		chat.NewMatchingService,
		chat.NewChannelService,

		chat.NewMelodyMatchConn,
		chat.NewMelodyChatConn,

		chat.NewRouter,
	)
	return &chat.Router{}, nil
}

func InitializeUploaderRouter() (*uploader.Router, error) {
	wire.Build(
		config.NewConfig,
		uploader.NewGinServer,
		uploader.NewRouter,
	)
	return &uploader.Router{}, nil
}
