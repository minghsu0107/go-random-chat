package common

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HttpLogrus struct {
	*log.Entry
}
type GrpcLogrus struct {
	*log.Entry
}

func InitLogging() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

func NewHttpLogrus() HttpLogrus {
	return HttpLogrus{log.WithField("protocol", "http")}
}

func NewGrpcLogrus() GrpcLogrus {
	return GrpcLogrus{log.WithField("protocol", "grpc")}
}
