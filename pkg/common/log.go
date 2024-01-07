package common

import (
	"io"
	"os"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/minghsu0107/go-random-chat/pkg/config"
)

type HttpLog struct {
	*slog.Logger
}
type GrpcLog struct {
	*slog.Logger
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: false,
	})))
}

func NewHttpLog(config *config.Config) (HttpLog, error) {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: false,
	}).WithAttrs([]slog.Attr{
		slog.String("proto", "http"),
	})
	logger := slog.New(logHandler)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)

	return HttpLog{logger}, nil
}

func NewGrpcLog(config *config.Config) (GrpcLog, error) {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: false,
	}).WithAttrs([]slog.Attr{
		slog.String("proto", "grpc"),
	})
	logger := slog.New(logHandler)

	return GrpcLog{logger}, nil
}
