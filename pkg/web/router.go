package web

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
)

type Router struct {
	obsInjector *common.ObservibilityInjector
	svr         *gin.Engine
	httpPort    string
	httpServer  *http.Server
}

func NewGinServer() *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware())

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: "web",
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewRouter(config *config.Config, obsInjector *common.ObservibilityInjector, svr *gin.Engine) *Router {
	return &Router{
		obsInjector: obsInjector,
		svr:         svr,
		httpPort:    config.Web.Http.Port,
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
}

func (r *Router) Run() {
	if err := r.obsInjector.Register("web"); err != nil {
		log.Error(err)
	}
	go func() {
		r.RegisterRoutes()
		addr := ":" + r.httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(r.svr, "web_http"),
		}
		log.Infoln("http server listening on ", addr)
		err := r.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}
func (r *Router) GracefulStop(ctx context.Context, done chan bool) {
	err := r.httpServer.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}

	log.Info("gracefully shutdowned")
	done <- true
}
