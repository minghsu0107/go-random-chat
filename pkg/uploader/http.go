package uploader

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

type HttpServer struct {
	name       string
	logger     common.HttpLogrus
	svr        *gin.Engine
	s3Endpoint string
	s3Bucket   string
	maxMemory  int64
	uploader   *s3manager.Uploader
	httpPort   string
	httpServer *http.Server
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

func NewHttpServer(name string, logger common.HttpLogrus, config *config.Config, svr *gin.Engine) common.HttpServer {
	initJWT(config)

	s3Endpoint := config.Uploader.S3.Endpoint
	s3Bucket := config.Uploader.S3.Bucket
	disableSSL := config.Uploader.S3.DisableSSL
	creds := credentials.NewStaticCredentials(config.Uploader.S3.AccessKey, config.Uploader.S3.SecretKey, "")

	awsConfig := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(s3Endpoint),
		Region:           aws.String(config.Uploader.S3.Region),
		DisableSSL:       aws.Bool(disableSSL),
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(3),
	}

	sess := session.Must(session.NewSession(awsConfig))
	return &HttpServer{
		name:       name,
		logger:     logger,
		svr:        svr,
		s3Endpoint: s3Endpoint,
		s3Bucket:   s3Bucket,
		maxMemory:  config.Uploader.Http.Server.MaxMemoryByte,
		uploader:   s3manager.NewUploader(sess),
		httpPort:   config.Uploader.Http.Server.Port,
	}
}

func initJWT(config *config.Config) {
	common.JwtSecret = config.Uploader.JWT.Secret
}

func (r *HttpServer) RegisterRoutes() {
	uploadGroup := r.svr.Group("/api/file")
	uploadGroup.Use(common.JWTAuth())
	{
		uploadGroup.POST("", r.UploadFile)
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
