package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

var (
	userPrefix    = "rc:user"
	sessionPrefix = "rc:session"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
	GetUserByOAuthEmail(ctx context.Context, authType AuthType, email string) (*User, error)
	SetUserSession(ctx context.Context, uid uint64, sid string) error
	GetUserIDBySession(ctx context.Context, sid string) (uint64, error)
}

type UserRepoImpl struct {
	r infra.RedisCache
}

func NewUserRepoImpl(r infra.RedisCache) *UserRepoImpl {
	return &UserRepoImpl{r}
}

func (repo *UserRepoImpl) CreateUser(ctx context.Context, user *User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if err = repo.r.Set(ctx, constructKey(userPrefix, user.ID), data); err != nil {
		return err
	}
	if user.AuthType != LocalAuth {
		if err = repo.r.Set(ctx, constructOAuthKey(user.AuthType, user.Email), data); err != nil {
			return err
		}
	}
	return nil
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

func (repo *UserRepoImpl) GetUserByOAuthEmail(ctx context.Context, authType AuthType, email string) (*User, error) {
	key := constructOAuthKey(authType, email)
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

func (repo *UserRepoImpl) SetUserSession(ctx context.Context, uid uint64, sid string) error {
	key := common.Join(sessionPrefix, ":", sid)
	return repo.r.Set(ctx, key, uid)
}

func (repo *UserRepoImpl) GetUserIDBySession(ctx context.Context, sid string) (uint64, error) {
	key := common.Join(sessionPrefix, ":", sid)
	var userID uint64
	exist, err := repo.r.Get(ctx, key, &userID)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, ErrSessionNotFound
	}
	return userID, nil
}

func constructKey(prefix string, id uint64) string {
	return common.Join(prefix, ":", strconv.FormatUint(id, 10))
}

func constructOAuthKey(authType AuthType, email string) string {
	return common.Join(userPrefix, ":", string(authType), ":", email)
}
