package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	svr        *gin.Engine
	httpPort   string
	httpServer *http.Server
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewGinServer() *gin.Engine {
	return gin.Default()
}

func NewRouter(config *config.Config, svr *gin.Engine) *Router {
	return &Router{
		svr:      svr,
		httpPort: config.Web.Http.Port,
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
	go func() {
		r.RegisterRoutes()
		addr := ":" + r.httpPort
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
}
func (r *Router) GracefulStop(ctx context.Context, done chan bool) {
	err := r.httpServer.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}

	log.Info("gracefully shutdowned")
	done <- true
}
