package upload

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	log "github.com/sirupsen/logrus"
)

var (
	httpPort = common.Getenv("HTTP_PORT", "5001")
)

type Router struct {
	svr        *gin.Engine
	httpServer *http.Server
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewRouter(svr *gin.Engine) *Router {
	return &Router{
		svr: svr,
	}
}

func (r *Router) RegisterRoutes() {
	uploadGroup := r.svr.Group("/api/file")
	uploadGroup.Use(common.JWTAuth())
	{
		uploadGroup.POST("", r.UploadFile)
	}
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
}
func (r *Router) GracefulStop(ctx context.Context, done chan bool) {
	err := r.httpServer.Shutdown(ctx)
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
