package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type HttpServer interface {
	RegisterRoutes()
	Run()
	GracefulStop(ctx context.Context) error
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
	obsInjector *ObservibilityInjector
}

func NewServer(name string, router Router, infraCloser InfraCloser, obsInjector *ObservibilityInjector) *Server {
	return &Server{name, router, infraCloser, obsInjector}
}

func (s *Server) Serve() {
	if err := s.obsInjector.Register(s.name); err != nil {
		log.Fatal(err)
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
		log.Error(err)
	}

	if TracerProvider != nil {
		err = TracerProvider.Shutdown(ctx)
		if err != nil {
			log.Error(err)
		}
	}

	if err = s.infraCloser.Close(); err != nil {
		log.Error(err)
	}

	log.Info("gracefully shutdowned")
	done <- true
}
