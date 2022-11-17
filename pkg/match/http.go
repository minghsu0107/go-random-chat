package match

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	"gopkg.in/olahol/melody.v1"

	doc "github.com/minghsu0107/go-random-chat/docs/match"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	sessUidKey = "sessuid"

	MelodyMatch MelodyMatchConn
)

type MelodyMatchConn struct {
	*melody.Melody
}

type HttpServer struct {
	name            string
	logger          common.HttpLogrus
	svr             *gin.Engine
	mm              MelodyMatchConn
	httpPort        string
	httpServer      *http.Server
	matchSubscriber *MatchSubscriber
	userSvc         UserService
	matchSvc        MatchingService
	serveSwag       bool
}

func NewMelodyMatchConn() MelodyMatchConn {
	MelodyMatch = MelodyMatchConn{
		melody.New(),
	}
	return MelodyMatch
}

func NewGinServer(name string, logger common.HttpLogrus, config *config.Config) *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.CorsMiddleware())
	svr.Use(common.LoggingMiddleware(logger))
	svr.Use(common.MaxAllowed(config.Match.Http.Server.MaxConn))

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: name,
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewHttpServer(name string, logger common.HttpLogrus, config *config.Config, svr *gin.Engine, mm MelodyMatchConn, matchSubscriber *MatchSubscriber, userSvc UserService, matchSvc MatchingService) common.HttpServer {
	initJWT(config)

	return &HttpServer{
		name:            name,
		logger:          logger,
		svr:             svr,
		mm:              mm,
		httpPort:        config.Match.Http.Server.Port,
		matchSubscriber: matchSubscriber,
		userSvc:         userSvc,
		matchSvc:        matchSvc,
		serveSwag:       config.Match.Http.Server.Swag,
	}
}

func (r *HttpServer) CookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := common.GetCookie(c, common.SessionIdCookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := r.userSvc.GetUserIDBySession(c.Request.Context(), sid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), common.UserKey, userID))
		c.Next()
	}
}

func initJWT(config *config.Config) {
	common.JwtSecret = config.Match.JWT.Secret
	common.JwtExpirationSecond = config.Match.JWT.ExpirationSecond
}

// @title           Match Service Swagger API
// @version         2.0
// @description     Match service API

// @contact.name   Ming Hsu
// @contact.email  minghsu0107@gmail.com

// @BasePath  /api
func (r *HttpServer) RegisterRoutes() {
	matchGroup := r.svr.Group("/api/match")
	{
		cookieAuthGroup := matchGroup.Group("")
		cookieAuthGroup.Use(r.CookieAuth())
		cookieAuthGroup.GET("", r.Match)

		forwardAuthGroup := matchGroup.Group("/forwardauth")
		forwardAuthGroup.Use(common.JWTAuth())
		{
			forwardAuthGroup.Any("", r.ForwardAuth)
		}
	}

	r.mm.HandleConnect(r.HandleMatchOnConnect)
	r.mm.HandleClose(r.HandleMatchOnClose)

	if r.serveSwag {
		matchGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName(doc.SwaggerInfomatch.InfoInstanceName)))
	}
}

func (r *HttpServer) Run() {
	go func() {
		addr := ":" + r.httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(r.svr, r.name+"_http"),
		}
		r.logger.Infoln("http server listening on ", addr)
		err := r.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			r.logger.Fatal(err)
		}
	}()
	go func() {
		err := r.matchSubscriber.Run()
		if err != nil {
			r.logger.Fatal(err)
		}
	}()
}
func (r *HttpServer) GracefulStop(ctx context.Context) error {
	err := r.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	err = r.matchSubscriber.GracefulStop()
	if err != nil {
		return err
	}
	return nil
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
