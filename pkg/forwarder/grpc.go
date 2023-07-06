package forwarder

import (
	"net"

	"google.golang.org/grpc"

	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	forwarderpb "github.com/minghsu0107/go-random-chat/proto/forwarder"
)

type GrpcServer struct {
	grpcPort      string
	logger        common.GrpcLogrus
	s             *grpc.Server
	forwardSvc    ForwardService
	msgSubscriber *MessageSubscriber
}

func NewGrpcServer(name string, logger common.GrpcLogrus, config *config.Config, forwardSvc ForwardService, msgSubscriber *MessageSubscriber) *GrpcServer {
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
		srv.logger.Infoln("grpc server listening on  ", addr)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Fatal(err)
		}
		if err := srv.s.Serve(lis); err != nil {
			srv.logger.Fatal(err)
		}
	}()
	go func() {
		err := srv.msgSubscriber.Run()
		if err != nil {
			srv.logger.Fatal(err)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GracefulStop()
	return srv.msgSubscriber.GracefulStop()
}
