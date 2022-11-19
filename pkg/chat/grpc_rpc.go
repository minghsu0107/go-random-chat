package chat

import (
	"context"

	chatpb "github.com/minghsu0107/go-random-chat/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *GrpcServer) CreateChannel(ctx context.Context, req *chatpb.CreateChannelRequest) (*chatpb.CreateChannelResponse, error) {
	channel, err := srv.chanSvc.CreateChannel(ctx)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &chatpb.CreateChannelResponse{
		ChannelId:   channel.ID,
		AccessToken: channel.AccessToken,
	}, nil
}

func (srv *GrpcServer) AddUserToChannel(ctx context.Context, req *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error) {
	if err := srv.userSvc.AddUserToChannel(ctx, req.ChannelId, req.UserId); err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &chatpb.AddUserResponse{}, nil
}
