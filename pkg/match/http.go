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
	matchSubscriber MatchSubscriber
	userSvc         UserService
	matchSvc        MatchingService
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

func NewHttpServer(name string, logger common.HttpLogrus, config *config.Config, svr *gin.Engine, mm MelodyMatchConn, matchSubscriber MatchSubscriber, userSvc UserService, matchSvc MatchingService) common.HttpServer {
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
	}
}

func initJWT(config *config.Config) {
	common.JwtSecret = config.Match.JWT.Secret
	common.JwtExpirationSecond = config.Match.JWT.ExpirationSecond
}

func (r *HttpServer) RegisterRoutes() {
	r.svr.GET("/api/match", r.Match)

	r.mm.HandleConnect(r.HandleMatchOnConnect)
	r.mm.HandleClose(r.HandleMatchOnClose)
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
		err := r.matchSubscriber.Subscribe()
		if err != nil {
			r.logger.Fatal(err)
		}
	}()
}
func (r *HttpServer) GracefulStop(ctx context.Context) error {
	r.matchSubscriber.Close()
	return r.httpServer.Shutdown(ctx)
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
