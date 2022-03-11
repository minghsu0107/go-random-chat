//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"github.com/minghsu0107/go-random-chat/pkg/uploader"
	"github.com/minghsu0107/go-random-chat/pkg/web"
)

func InitializeWebServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		web.NewGinServer,
		web.NewHttpServer,
		web.NewRouter,
		web.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeChatServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,

		infra.NewRedisClient,
		infra.NewRedisCache,

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

		chat.NewGinServer,
		chat.NewHttpServer,
		chat.NewRouter,
		chat.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeUploaderServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		uploader.NewGinServer,
		uploader.NewHttpServer,
		uploader.NewRouter,
		uploader.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}
