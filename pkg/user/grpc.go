package user

import (
	"net"

	"google.golang.org/grpc"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	userpb "github.com/minghsu0107/go-random-chat/proto/user"
)

type GrpcServer struct {
	grpcPort string
	logger   common.GrpcLogrus
	s        *grpc.Server
	userSvc  UserService
}

func NewGrpcServer(logger common.GrpcLogrus, config *config.Config, userSvc UserService) common.GrpcServer {
	srv := &GrpcServer{
		grpcPort: config.User.Grpc.Server.Port,
		logger:   logger,
		userSvc:  userSvc,
	}
	srv.s = transport.InitializeGrpcServer(srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	userpb.RegisterUserServiceServer(srv.s, srv)
	grpc_prometheus.Register(srv.s)
}

func (srv *GrpcServer) Run() {
	go func() {
		addr := "0.0.0.0:" + srv.grpcPort
		srv.logger.Infoln("grpc server listening on  ", addr)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Fatal(err)
		}
		if err := srv.s.Serve(lis); err != nil {
			srv.logger.Fatal(err)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GracefulStop()
	return nil
}
