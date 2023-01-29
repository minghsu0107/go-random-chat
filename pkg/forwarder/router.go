package forwarder

import (
	"context"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type Router struct {
	grpcServer common.GrpcServer
}

func NewRouter(grpcServer common.GrpcServer) common.Router {
	return &Router{grpcServer}
}

func (r *Router) Run() {
	r.grpcServer.Register()
	r.grpcServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	r.grpcServer.GracefulStop()
	return nil
}
