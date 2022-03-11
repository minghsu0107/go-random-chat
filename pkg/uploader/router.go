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
	log "github.com/sirupsen/logrus"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

type Router struct {
	obsInjector *common.ObservibilityInjector
	svr         *gin.Engine
	s3Endpoint  string
	s3Bucket    string
	uploader    *s3manager.Uploader
	httpPort    string
	httpServer  *http.Server
}

func NewGinServer() *gin.Engine {
	common.InitLogging()

	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(common.LoggingMiddleware())
	svr.Use(common.CORSMiddleware())

	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: "uploader",
		}),
	})
	svr.Use(ginmiddleware.Handler("", mdlw))
	return svr
}

func NewRouter(config *config.Config, obsInjector *common.ObservibilityInjector, svr *gin.Engine) *Router {
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
	return &Router{
		obsInjector: obsInjector,
		svr:         svr,
		s3Endpoint:  s3Endpoint,
		s3Bucket:    s3Bucket,
		uploader:    s3manager.NewUploader(sess),
		httpPort:    config.Uploader.Http.Port,
	}
}

func initJWT(config *config.Config) {
	common.JwtSecret = config.Uploader.JWT.Secret
}

func (r *Router) RegisterRoutes() {
	uploadGroup := r.svr.Group("/api/file")
	uploadGroup.Use(common.JWTAuth())
	{
		uploadGroup.POST("", r.UploadFile)
	}
}

func (r *Router) Run() {
	if err := r.obsInjector.Register("uploader"); err != nil {
		log.Error(err)
	}
	go func() {
		r.RegisterRoutes()
		addr := ":" + r.httpPort
		r.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(r.svr, "uploader_http"),
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
