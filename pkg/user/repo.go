package user

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

var (
	userPrefix = "rc:user"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
}

type UserRepoImpl struct {
	r infra.RedisCache
}

func NewUserRepo(r infra.RedisCache) UserRepo {
	return &UserRepoImpl{r}
}
func (repo *UserRepoImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = repo.r.Set(ctx, constructKey(userPrefix, user.ID), data)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
func (repo *UserRepoImpl) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	key := constructKey(userPrefix, userID)
	var user User
	exist, err := repo.r.Get(ctx, key, &user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func constructKey(prefix string, id uint64) string {
	return prefix + ":" + strconv.FormatUint(id, 10)
}
