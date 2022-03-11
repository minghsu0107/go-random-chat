//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/uploader"
	"github.com/minghsu0107/go-random-chat/pkg/web"
)

func InitializeWebRouter() (*web.Router, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		web.NewGinServer,
		web.NewRouter,
	)
	return &web.Router{}, nil
}

func InitializeChatRouter() (*chat.Router, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,

		chat.NewRedisClient,
		chat.NewRedisCache,

		chat.NewKafkaPublisher,
		chat.NewKafkaSubscriber,

		chat.NewRedisUserRepo,
		chat.NewMessageRepo,
		chat.NewRedisChannelRepo,
		chat.NewMatchingRepo,

		chat.NewMessageSubscriber,
		chat.NewMatchSubscriber,

		chat.NewSonyFlake,

		chat.NewUserService,
		chat.NewMessageService,
		chat.NewMatchingService,
		chat.NewChannelService,

		chat.NewMelodyMatchConn,
		chat.NewMelodyChatConn,

		chat.NewGinServer,
		chat.NewRouter,
	)
	return &chat.Router{}, nil
}

func InitializeUploaderRouter() (*uploader.Router, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		uploader.NewGinServer,
		uploader.NewRouter,
	)
	return &uploader.Router{}, nil
}
