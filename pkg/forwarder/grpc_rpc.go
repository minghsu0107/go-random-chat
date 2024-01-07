package forwarder

import (
	"context"

	forwarderpb "github.com/minghsu0107/go-random-chat/proto/forwarder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *GrpcServer) RegisterChannelSession(ctx context.Context, req *forwarderpb.RegisterChannelSessionRequest) (*forwarderpb.RegisterChannelSessionResponse, error) {
	if err := srv.forwardSvc.RegisterChannelSession(ctx, req.ChannelId, req.UserId, req.Subscriber); err != nil {
		srv.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &forwarderpb.RegisterChannelSessionResponse{}, nil
}

func (srv *GrpcServer) RemoveChannelSession(ctx context.Context, req *forwarderpb.RemoveChannelSessionRequest) (*forwarderpb.RemoveChannelSessionResponse, error) {
	if err := srv.forwardSvc.RemoveChannelSession(ctx, req.ChannelId, req.UserId); err != nil {
		srv.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &forwarderpb.RemoveChannelSessionResponse{}, nil
}
