package forwarder

import (
	"net"
	"os"

	"google.golang.org/grpc"

	"log/slog"

	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	forwarderpb "github.com/minghsu0107/go-random-chat/proto/forwarder"
)

type GrpcServer struct {
	grpcPort      string
	logger        common.GrpcLog
	s             *grpc.Server
	forwardSvc    ForwardService
	msgSubscriber *MessageSubscriber
}

func NewGrpcServer(name string, logger common.GrpcLog, config *config.Config, forwardSvc ForwardService, msgSubscriber *MessageSubscriber) *GrpcServer {
	srv := &GrpcServer{
		grpcPort:      config.Forwarder.Grpc.Server.Port,
		logger:        logger,
		forwardSvc:    forwardSvc,
		msgSubscriber: msgSubscriber,
	}
	srv.s = transport.InitializeGrpcServer(name, srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	srv.msgSubscriber.RegisterHandler()

	forwarderpb.RegisterForwardServiceServer(srv.s, srv)
}

func (srv *GrpcServer) Run() {
	go func() {
		addr := "0.0.0.0:" + srv.grpcPort
		srv.logger.Info("grpc server listening", slog.String("addr", addr))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
		if err := srv.s.Serve(lis); err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
	go func() {
		err := srv.msgSubscriber.Run()
		if err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GracefulStop()
	return srv.msgSubscriber.GracefulStop()
}
