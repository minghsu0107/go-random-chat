package user

import (
	"context"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type UserService interface {
	CreateUser(ctx context.Context, userName string) (*User, error)
	GetUser(ctx context.Context, uid uint64) (*User, error)
}

type UserServiceImpl struct {
	userRepo UserRepo
	sf       common.IDGenerator
}

func NewUserService(userRepo UserRepo, sf common.IDGenerator) UserService {
	return &UserServiceImpl{userRepo, sf}
}
func (svc *UserServiceImpl) CreateUser(ctx context.Context, userName string) (*User, error) {
	userID, err := svc.sf.NextID()
	if err != nil {
		return nil, err
	}
	return svc.userRepo.CreateUser(ctx, &User{
		ID:   userID,
		Name: userName,
	})
}
func (svc *UserServiceImpl) GetUser(ctx context.Context, uid uint64) (*User, error) {
	return svc.userRepo.GetUserByID(ctx, uid)
}
