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

	doc "github.com/minghsu0107/go-random-chat/docs/uploader"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	serveSwag  bool
}

func NewGinServer(name string, logger common.HttpLogrus, config *config.Config) *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware(logger))
	svr.Use(common.CORSMiddleware())
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
		serveSwag:  config.Uploader.Http.Server.Swag,
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
		fileGroup := uploaderGroup.Group("/file")
		fileGroup.Use(common.JWTForwardAuth())
		{
			fileGroup.POST("", r.UploadFile)
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
