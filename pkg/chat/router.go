package chat

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
)

var (
	httpPort        = getenv("HTTP_PORT", "5000")
	maxAllowedConns int64
	sessUidKey      = "sessuid"
	sessCidKey      = "sesscid"
	MelodyMatch     *melody.Melody
	MelodyChat      *melody.Melody
)

type Router struct {
	svr             *gin.Engine
	mm              *melody.Melody
	mc              *melody.Melody
	httpServer      *http.Server
	matchSubscriber MatchSubscriber
	msgSubscriber   MessageSubscriber
	userSvc         UserService
	msgSvc          MessageService
	matchSvc        MatchingService
	chanSvc         ChannelService
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	var err error
	maxAllowedConns, err = strconv.ParseInt(getenv("MAX_ALLOWED_CONNS", "200"), 10, 64)
	if err != nil {
		panic(err)
	}
}

func NewRouter(svr *gin.Engine, mm, mc *melody.Melody, matchSubscriber MatchSubscriber, msgSubscriber MessageSubscriber, userSvc UserService, msgSvc MessageService, matchSvc MatchingService, chanSvc ChannelService) *Router {
	return &Router{
		svr:             svr,
		mm:              mm,
		mc:              mc,
		matchSubscriber: matchSubscriber,
		msgSubscriber:   msgSubscriber,
		userSvc:         userSvc,
		msgSvc:          msgSvc,
		matchSvc:        matchSvc,
		chanSvc:         chanSvc,
	}
}

func (r *Router) RegisterRoutes() {
	r.svr.LoadHTMLGlob("web/html/*")
	r.svr.Static("/assets", "./web/assets")
	r.svr.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.svr.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	r.svr.GET("/api/match", r.Match)
	r.svr.GET("/api/chat", r.StartChat)

	userGroup := r.svr.Group("/api/user")
	{
		userGroup.POST("", r.CreateUser)
		userGroup.GET("/:uid/name", r.GetUserName)
	}
	usersGroup := r.svr.Group("/api/users")
	usersGroup.Use(JWTAuth())
	{
		usersGroup.GET("", r.GetChannelUsers)
		usersGroup.GET("/online", r.GetOnlineUsers)
	}
	channelGroup := r.svr.Group("/api/channel")
	channelGroup.Use(JWTAuth())
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
	go func() {
		r.RegisterRoutes()
		addr := ":" + httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: r.svr,
		}
		log.Infoln("http server listening on ", addr)
		err := r.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		err := r.matchSubscriber.Subscribe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		err := r.msgSubscriber.Subscribe()
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
	r.matchSubscriber.Close()
	r.msgSubscriber.Close()
	if err = RedisClient.Close(); err != nil {
		log.Error(err)
	}

	log.Info("gracefully shutdowned")
	done <- true
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, ErrResponse{
		Message: message,
	})
}
