package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

type HttpServer struct {
	name       string
	logger     common.HttpLog
	svr        *gin.Engine
	httpPort   string
	httpServer *http.Server
}

func NewGinServer(name string, logger common.HttpLog) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware(logger))

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: name,
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewHttpServer(name string, logger common.HttpLog, config *config.Config, svr *gin.Engine) *HttpServer {
	return &HttpServer{
		name:     name,
		logger:   logger,
		svr:      svr,
		httpPort: config.Web.Http.Server.Port,
	}
}

func (r *HttpServer) RegisterRoutes() {
	r.svr.LoadHTMLGlob("web/html/*")
	r.svr.Static("/assets", "./web/assets")
	r.svr.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.svr.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})
}

func (r *HttpServer) Run() {
	go func() {
		addr := ":" + r.httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(r.svr, r.name+"_http"),
		}
		r.logger.Info("http server listening", slog.String("addr", addr))
		err := r.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			r.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}
func (r *HttpServer) GracefulStop(ctx context.Context) error {
	return r.httpServer.Shutdown(ctx)
}
