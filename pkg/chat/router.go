package chat

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	log "github.com/sirupsen/logrus"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	"gopkg.in/olahol/melody.v1"
)

var (
	sessUidKey = "sessuid"
	sessCidKey = "sesscid"

	MelodyMatch MelodyMatchConn
	MelodyChat  MelodyChatConn
)

type MelodyMatchConn struct {
	*melody.Melody
}
type MelodyChatConn struct {
	*melody.Melody
}

type Router struct {
	obsInjector     *common.ObservibilityInjector
	svr             *gin.Engine
	mm              MelodyMatchConn
	mc              MelodyChatConn
	httpPort        string
	httpServer      *http.Server
	matchSubscriber *MatchSubscriber
	msgSubscriber   *MessageSubscriber
	userSvc         UserService
	msgSvc          MessageService
	matchSvc        MatchingService
	chanSvc         ChannelService
}

func NewMelodyMatchConn() MelodyMatchConn {
	MelodyMatch = MelodyMatchConn{
		melody.New(),
	}
	return MelodyMatch
}

func NewMelodyChatConn(config *config.Config) MelodyChatConn {
	m := melody.New()
	m.Config.MaxMessageSize = config.Chat.Message.MaxSizeByte
	MelodyChat = MelodyChatConn{
		m,
	}
	return MelodyChat
}

func NewGinServer(config *config.Config) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware())
	svr.Use(common.MaxAllowed(config.Chat.Http.MaxConn))
	svr.Use(common.CORSMiddleware())

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: "chat",
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewRouter(config *config.Config, obsInjector *common.ObservibilityInjector, svr *gin.Engine, mm MelodyMatchConn, mc MelodyChatConn, matchSubscriber *MatchSubscriber, msgSubscriber *MessageSubscriber, userSvc UserService, msgSvc MessageService, matchSvc MatchingService, chanSvc ChannelService) *Router {
	common.InitLogging()
	initJWT(config)

	return &Router{
		obsInjector:     obsInjector,
		svr:             svr,
		mm:              mm,
		mc:              mc,
		httpPort:        config.Chat.Http.Port,
		matchSubscriber: matchSubscriber,
		msgSubscriber:   msgSubscriber,
		userSvc:         userSvc,
		msgSvc:          msgSvc,
		matchSvc:        matchSvc,
		chanSvc:         chanSvc,
	}
}

func initJWT(config *config.Config) {
	common.JwtSecret = config.Chat.JWT.Secret
	common.JwtExpirationSecond = config.Chat.JWT.ExpirationSecond
}

func (r *Router) RegisterRoutes() {
	r.svr.GET("/api/match", r.Match)
	r.svr.GET("/api/chat", r.StartChat)

	userGroup := r.svr.Group("/api/user")
	{
		userGroup.POST("", r.CreateUser)
		userGroup.GET("/:uid/name", r.GetUserName)
	}
	usersGroup := r.svr.Group("/api/users")
	usersGroup.Use(common.JWTAuth())
	{
		usersGroup.GET("", r.GetChannelUsers)
		usersGroup.GET("/online", r.GetOnlineUsers)
	}
	channelGroup := r.svr.Group("/api/channel")
	channelGroup.Use(common.JWTAuth())
	{
		channelGroup.GET("/messages", r.ListMessages)
		channelGroup.DELETE("", r.DeleteChannel)
	}

	r.mm.HandleConnect(r.HandleMatchOnConnect)
	r.mm.HandleClose(r.HandleMatchOnClose)

	r.mc.HandleMessage(r.HandleChatOnMessage)
	r.mc.HandleConnect(r.HandleChatOnConnect)
	r.mc.HandleClose(r.HandleChatOnClose)
}

func (r *Router) Run() {
	if err := r.obsInjector.Register("chat"); err != nil {
		log.Error(err)
	}
	go func() {
		r.RegisterRoutes()
		addr := ":" + r.httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(r.svr, "chat_http"),
		}
		log.Infoln("http server listening on ", addr)
		err := r.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		err := r.matchSubscriber.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := r.msgSubscriber.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
func (r *Router) GracefulStop(ctx context.Context, done chan bool) {
	err := r.httpServer.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
	err = r.matchSubscriber.GracefulStop()
	if err != nil {
		log.Error(err)
	}
	err = r.msgSubscriber.GracefulStop()
	if err != nil {
		log.Error(err)
	}
	err = RedisClient.Close()
	if err != nil {
		log.Error(err)
	}

	log.Info("gracefully shutdowned")
	done <- true
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
