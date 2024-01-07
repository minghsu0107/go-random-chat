package common

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer interface {
	RegisterRoutes()
	Run()
	GracefulStop(ctx context.Context) error
}

type GrpcServer interface {
	Register()
	Run()
	GracefulStop() error
}

type Router interface {
	Run()
	GracefulStop(ctx context.Context) error
}

type InfraCloser interface {
	Close() error
}

type Server struct {
	name        string
	router      Router
	infraCloser InfraCloser
	obsInjector *ObservabilityInjector
}

func NewServer(name string, router Router, infraCloser InfraCloser, obsInjector *ObservabilityInjector) *Server {
	return &Server{name, router, infraCloser, obsInjector}
}

func (s *Server) Serve() {
	if err := s.obsInjector.Register(s.name); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	s.router.Run()

	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.GracefulStop(ctx, done)
	}()

	<-done
}

func (s *Server) GracefulStop(ctx context.Context, done chan bool) {
	err := s.router.GracefulStop(ctx)
	if err != nil {
		slog.Error(err.Error())
	}

	if TracerProvider != nil {
		err = TracerProvider.Shutdown(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	if err = s.infraCloser.Close(); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("gracefully shutdowned")
	done <- true
}
