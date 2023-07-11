package common

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	log "github.com/sirupsen/logrus"
)

type HttpLogrus struct {
	*log.Entry
}
type GrpcLogrus struct {
	*log.Entry
}

func NewHttpLogrus(config *config.Config) (HttpLogrus, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)

	if err := initLogging(config); err != nil {
		return HttpLogrus{}, err
	}

	return HttpLogrus{log.WithField("protocol", "http")}, nil
}

func NewGrpcLogrus(config *config.Config) (GrpcLogrus, error) {
	if err := initLogging(config); err != nil {
		return GrpcLogrus{}, err
	}

	return GrpcLogrus{log.WithField("protocol", "grpc")}, nil
}

func initLogging(config *config.Config) error {
	logrusLevel, err := log.ParseLevel(config.Logging.Level)
	if err != nil {
		return err
	}
	log.SetOutput(os.Stderr)
	log.SetLevel(logrusLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	return nil
}
