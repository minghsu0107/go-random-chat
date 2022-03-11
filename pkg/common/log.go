package common

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitLogging() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}
