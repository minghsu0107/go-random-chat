package uploader

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/redis/go-redis/v9"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"

	doc "github.com/minghsu0107/go-random-chat/docs/uploader"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ChannelUploadRateLimiter struct {
	*common.RateLimiter
}

func NewChannelUploadRateLimiter(rc redis.UniversalClient, config *config.Config) ChannelUploadRateLimiter {
	return ChannelUploadRateLimiter{
		common.NewRateLimiter(
			rc,
			config.Uploader.RateLimit.ChannelUpload.Rps,
			config.Uploader.RateLimit.ChannelUpload.Burst,
			time.Duration(config.Redis.ExpirationHour)*time.Hour,
		),
	}
}

type HttpServer struct {
	name                     string
	logger                   common.HttpLogrus
	svr                      *gin.Engine
	s3Endpoint               string
	s3Bucket                 string
	maxMemory                int64
	uploader                 *manager.Uploader
	httpPort                 string
	httpServer               *http.Server
	channelUploadRateLimiter ChannelUploadRateLimiter
	serveSwag                bool
}

func NewGinServer(name string, logger common.HttpLogrus, config *config.Config) *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.CorsMiddleware())
	svr.Use(common.LoggingMiddleware(logger))
	svr.Use(common.LimitBodySize(config.Uploader.Http.Server.MaxBodyByte))

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: name,
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewHttpServer(name string, logger common.HttpLogrus, config *config.Config, svr *gin.Engine, channelUploadRateLimiter ChannelUploadRateLimiter) *HttpServer {
	s3Endpoint := config.Uploader.S3.Endpoint
	s3Bucket := config.Uploader.S3.Bucket
	creds := credentials.NewStaticCredentialsProvider(config.Uploader.S3.AccessKey, config.Uploader.S3.SecretKey, "")
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               s3Endpoint,
			SigningRegion:     config.Uploader.S3.Region,
			HostnameImmutable: true,
		}, nil
	})
	awsConfig := aws.Config{
		Credentials:                 creds,
		EndpointResolverWithOptions: customResolver,
		Region:                      config.Uploader.S3.Region,
		RetryMaxAttempts:            3,
	}

	return &HttpServer{
		name:                     name,
		logger:                   logger,
		svr:                      svr,
		s3Endpoint:               s3Endpoint,
		s3Bucket:                 s3Bucket,
		maxMemory:                config.Uploader.Http.Server.MaxMemoryByte,
		uploader:                 manager.NewUploader(s3.NewFromConfig(awsConfig)),
		httpPort:                 config.Uploader.Http.Server.Port,
		channelUploadRateLimiter: channelUploadRateLimiter,
		serveSwag:                config.Uploader.Http.Server.Swag,
	}
}

func (r *HttpServer) ChannelUploadRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		channelID, ok := c.Request.Context().Value(common.ChannelKey).(uint64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		allow, err := r.channelUploadRateLimiter.Allow(c.Request.Context(), strconv.FormatUint(channelID, 10))
		if err != nil {
			r.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !allow {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

// @title           Uploader Service Swagger API
// @version         2.0
// @description     Uploader service API

// @contact.name   Ming Hsu
// @contact.email  minghsu0107@gmail.com

// @BasePath  /api
func (r *HttpServer) RegisterRoutes() {
	uploaderGroup := r.svr.Group("/api/uploader")
	{
		uploadFileGroup := uploaderGroup.Group("/files")
		uploadFileGroup.Use(common.JWTForwardAuth())
		uploadFileGroup.Use(r.ChannelUploadRateLimit())
		{
			uploadFileGroup.POST("", r.UploadFiles)
		}
	}
	if r.serveSwag {
		uploaderGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName(doc.SwaggerInfouploader.InfoInstanceName)))
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
