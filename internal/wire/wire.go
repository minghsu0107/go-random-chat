//+build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"github.com/minghsu0107/go-random-chat/pkg/match"
	"github.com/minghsu0107/go-random-chat/pkg/uploader"
	"github.com/minghsu0107/go-random-chat/pkg/web"
)

func InitializeWebServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		common.NewHttpLogrus,
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
		common.NewHttpLogrus,
		common.NewGrpcLogrus,

		infra.NewRedisClient,
		infra.NewRedisCache,

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,

		chat.NewUserRepo,
		chat.NewMessageRepo,
		chat.NewChannelRepo,

		chat.NewMessageSubscriber,

		chat.NewSonyFlake,

		chat.NewUserService,
		chat.NewMessageService,
		chat.NewChannelService,

		chat.NewMelodyChatConn,

		chat.NewGinServer,
		chat.NewHttpServer,
		chat.NewGrpcServer,
		chat.NewRouter,
		chat.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeMatchServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		common.NewHttpLogrus,

		infra.NewRedisClient,
		infra.NewRedisCache,

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,

		match.NewChatClientConn,

		match.NewChannelRepo,
		match.NewUserRepo,
		match.NewMatchingRepo,

		match.NewMatchSubscriber,

		match.NewUserService,
		match.NewMatchingService,

		match.NewMelodyMatchConn,

		match.NewGinServer,
		match.NewHttpServer,
		match.NewRouter,
		match.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeUploaderServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservibilityInjector,
		common.NewHttpLogrus,
		uploader.NewGinServer,
		uploader.NewHttpServer,
		uploader.NewRouter,
		uploader.NewInfraCloser,
		common.NewServer,
	)
	return &common.Server{}, nil
}
