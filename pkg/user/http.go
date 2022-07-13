package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"

	doc "github.com/minghsu0107/go-random-chat/docs/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	name       string
	logger     common.HttpLogrus
	svr        *gin.Engine
	httpPort   string
	httpServer *http.Server
	userSvc    UserService
	serveSwag  bool
}

func NewGinServer(name string, logger common.HttpLogrus, config *config.Config) *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware(logger))
	svr.Use(common.CORSMiddleware())

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: name,
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewHttpServer(name string, logger common.HttpLogrus, config *config.Config, svr *gin.Engine, userSvc UserService) common.HttpServer {
	return &HttpServer{
		name:      name,
		logger:    logger,
		svr:       svr,
		httpPort:  config.User.Http.Server.Port,
		userSvc:   userSvc,
		serveSwag: config.User.Http.Server.Swag,
	}
}

// @title           User Service Swagger API
// @version         2.0
// @description     User service API

// @contact.name   Ming Hsu
// @contact.email  minghsu0107@gmail.com

// @BasePath  /api
func (r *HttpServer) RegisterRoutes() {
	userGroup := r.svr.Group("/api/user")
	{
		userGroup.POST("", r.CreateUser)
		userGroup.GET("/:uid/name", r.GetUserName)
	}
	if r.serveSwag {
		userGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName(doc.SwaggerInfouser.InfoInstanceName)))
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
}
func (r *HttpServer) GracefulStop(ctx context.Context) error {
	return r.httpServer.Shutdown(ctx)
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
