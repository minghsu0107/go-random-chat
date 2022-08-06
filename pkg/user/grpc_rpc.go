package user

import (
	"context"
	"errors"

	userpb "github.com/minghsu0107/go-random-chat/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *GrpcServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := srv.userSvc.GetUser(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return &userpb.GetUserResponse{
				Exist: false,
			}, nil
		}
		srv.logger.Error(err)
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return &userpb.GetUserResponse{
		Exist: true,
		User: &userpb.User{
			Id:   user.ID,
			Name: user.Name,
		},
	}, nil
}
