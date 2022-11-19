package user

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"

	doc "github.com/minghsu0107/go-random-chat/docs/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type HttpServer struct {
	name              string
	logger            common.HttpLogrus
	svr               *gin.Engine
	httpPort          string
	httpServer        *http.Server
	userSvc           UserService
	serveSwag         bool
	googleOauthConfig *oauth2.Config
	oauthCookieConfig config.CookieConfig
	authCookieConfig  config.CookieConfig
}

func NewGinServer(name string, logger common.HttpLogrus, config *config.Config) *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.CorsMiddleware())
	svr.Use(common.LoggingMiddleware(logger))

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
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  config.User.OAuth.Google.RedirectUrl,
			ClientID:     config.User.OAuth.Google.ClientId,
			ClientSecret: config.User.OAuth.Google.ClientSecret,
			Scopes:       strings.Split(config.User.OAuth.Google.Scopes, ","),
			Endpoint:     google.Endpoint,
		},
		oauthCookieConfig: config.User.OAuth.Cookie,
		authCookieConfig:  config.User.Auth.Cookie,
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

// @title           User Service Swagger API
// @version         2.0
// @description     User service API

// @contact.name   Ming Hsu
// @contact.email  minghsu0107@gmail.com

// @BasePath  /api
func (r *HttpServer) RegisterRoutes() {
	userGroup := r.svr.Group("/api/user")
	{
		userGroup.POST("", r.CreateLocalUser)

		cookieAuthGroup := userGroup.Group("")
		cookieAuthGroup.Use(r.CookieAuth())
		cookieAuthGroup.GET("", r.GetUser)
		cookieAuthGroup.GET("/me", r.GetUserMe)

		userGroup.GET("/oauth2/google/login", r.OAuthGoogleLogin)
		userGroup.GET("/oauth2/google/callback", r.OAuthGoogleCallback)
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
