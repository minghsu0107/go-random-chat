//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/google/wire"
	"github.com/minghsu0107/go-random-chat/pkg/chat"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/forwarder"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"github.com/minghsu0107/go-random-chat/pkg/match"
	"github.com/minghsu0107/go-random-chat/pkg/uploader"
	"github.com/minghsu0107/go-random-chat/pkg/user"
	"github.com/minghsu0107/go-random-chat/pkg/web"
)

func InitializeWebServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewHttpLog,

		web.NewGinServer,

		web.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*web.HttpServer)),
		web.NewRouter,
		wire.Bind(new(common.Router), new(*web.Router)),
		web.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*web.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeChatServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewHttpLog,
		common.NewGrpcLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,
		infra.NewBrokerRouter,

		infra.NewCassandraSession,

		chat.NewUserClientConn,
		chat.NewForwarderClientConn,

		chat.NewUserRepoImpl,
		wire.Bind(new(chat.UserRepo), new(*chat.UserRepoImpl)),
		chat.NewMessageRepoImpl,
		wire.Bind(new(chat.MessageRepo), new(*chat.MessageRepoImpl)),
		chat.NewChannelRepoImpl,
		wire.Bind(new(chat.ChannelRepo), new(*chat.ChannelRepoImpl)),
		chat.NewForwardRepoImpl,
		wire.Bind(new(chat.ForwardRepo), new(*chat.ForwardRepoImpl)),

		chat.NewUserRepoCacheImpl,
		wire.Bind(new(chat.UserRepoCache), new(*chat.UserRepoCacheImpl)),
		chat.NewMessageRepoCacheImpl,
		wire.Bind(new(chat.MessageRepoCache), new(*chat.MessageRepoCacheImpl)),
		chat.NewChannelRepoCacheImpl,
		wire.Bind(new(chat.ChannelRepoCache), new(*chat.ChannelRepoCacheImpl)),

		chat.NewMessageSubscriber,

		common.NewSonyFlake,

		chat.NewUserServiceImpl,
		wire.Bind(new(chat.UserService), new(*chat.UserServiceImpl)),
		chat.NewMessageServiceImpl,
		wire.Bind(new(chat.MessageService), new(*chat.MessageServiceImpl)),
		chat.NewChannelServiceImpl,
		wire.Bind(new(chat.ChannelService), new(*chat.ChannelServiceImpl)),
		chat.NewForwardServiceImpl,
		wire.Bind(new(chat.ForwardService), new(*chat.ForwardServiceImpl)),

		chat.NewMelodyChatConn,

		chat.NewGinServer,

		chat.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*chat.HttpServer)),
		chat.NewGrpcServer,
		wire.Bind(new(common.GrpcServer), new(*chat.GrpcServer)),
		chat.NewRouter,
		wire.Bind(new(common.Router), new(*chat.Router)),
		chat.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*chat.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeForwarderServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewGrpcLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,
		infra.NewBrokerRouter,

		forwarder.NewForwardRepoImpl,
		wire.Bind(new(forwarder.ForwardRepo), new(*forwarder.ForwardRepoImpl)),

		forwarder.NewForwardServiceImpl,
		wire.Bind(new(forwarder.ForwardService), new(*forwarder.ForwardServiceImpl)),

		forwarder.NewMessageSubscriber,

		forwarder.NewGrpcServer,
		wire.Bind(new(common.GrpcServer), new(*forwarder.GrpcServer)),
		forwarder.NewRouter,
		wire.Bind(new(common.Router), new(*forwarder.Router)),
		forwarder.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*forwarder.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeMatchServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewHttpLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,
		infra.NewBrokerRouter,

		match.NewChatClientConn,
		match.NewUserClientConn,

		match.NewChannelRepoImpl,
		wire.Bind(new(match.ChannelRepo), new(*match.ChannelRepoImpl)),
		match.NewUserRepoImpl,
		wire.Bind(new(match.UserRepo), new(*match.UserRepoImpl)),
		match.NewMatchingRepoImpl,
		wire.Bind(new(match.MatchingRepo), new(*match.MatchingRepoImpl)),

		match.NewMatchSubscriber,

		match.NewUserServiceImpl,
		wire.Bind(new(match.UserService), new(*match.UserServiceImpl)),
		match.NewMatchingServiceImpl,
		wire.Bind(new(match.MatchingService), new(*match.MatchingServiceImpl)),

		match.NewMelodyMatchConn,

		match.NewGinServer,

		match.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*match.HttpServer)),
		match.NewRouter,
		wire.Bind(new(common.Router), new(*match.Router)),
		match.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*match.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeUploaderServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewHttpLog,

		infra.NewRedisClient,

		uploader.NewGinServer,

		uploader.NewChannelUploadRateLimiter,

		uploader.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*uploader.HttpServer)),
		uploader.NewRouter,
		wire.Bind(new(common.Router), new(*uploader.Router)),
		uploader.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*uploader.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeUserServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewObservabilityInjector,
		common.NewHttpLog,
		common.NewGrpcLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		user.NewUserRepoImpl,
		wire.Bind(new(user.UserRepo), new(*user.UserRepoImpl)),

		common.NewSonyFlake,

		user.NewUserServiceImpl,
		wire.Bind(new(user.UserService), new(*user.UserServiceImpl)),

		user.NewGinServer,

		user.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*user.HttpServer)),
		user.NewGrpcServer,
		wire.Bind(new(common.GrpcServer), new(*user.GrpcServer)),
		user.NewRouter,
		wire.Bind(new(common.Router), new(*user.Router)),
		user.NewInfraCloser,
		wire.Bind(new(common.InfraCloser), new(*user.InfraCloser)),
		common.NewServer,
	)
	return &common.Server{}, nil
}
