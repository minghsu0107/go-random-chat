package user

import (
	"context"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type Router struct {
	httpServer common.HttpServer
	grpcServer common.GrpcServer
}

func NewRouter(httpServer common.HttpServer, grpcServer common.GrpcServer) common.Router {
	return &Router{httpServer, grpcServer}
}

func (r *Router) Run() {
	r.httpServer.RegisterRoutes()
	r.httpServer.Run()

	r.grpcServer.Register()
	r.grpcServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	r.grpcServer.GracefulStop()
	return r.httpServer.GracefulStop(ctx)
}
