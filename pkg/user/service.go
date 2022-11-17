package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type UserService interface {
	GetUserDataFromGoogle(ctx context.Context, code string) (*GoogleUserPresenter, error)
	GetOrCreateGoogleUser(ctx context.Context, email, userName string) (*User, error)
	CreateUser(ctx context.Context, userName string) (*User, error)
	SetUserSession(ctx context.Context, uid uint64) (string, error)
	GetUserByID(ctx context.Context, uid uint64) (*User, error)
	GetUserIDBySession(ctx context.Context, sid string) (uint64, error)
}

type UserServiceImpl struct {
	userRepo UserRepo
	sf       common.IDGenerator
}

func NewUserService(userRepo UserRepo, sf common.IDGenerator) UserService {
	return &UserServiceImpl{userRepo, sf}
}

func (svc *UserServiceImpl) GetUserDataFromGoogle(ctx context.Context, accessToken string) (*GoogleUserPresenter, error) {
	req, err := http.NewRequest("GET", common.Join(oauthGoogleUrlAPI, accessToken), nil)
	if err != nil {
		return nil, fmt.Errorf("create http request error: %w", err)
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read google user response: %w", err)
	}
	var googleUser GoogleUserPresenter
	if err := json.Unmarshal(contents, &googleUser); err != nil {
		return nil, fmt.Errorf("failed marshal google user response: %w", err)
	}
	return &googleUser, nil
}

func (svc *UserServiceImpl) CreateUser(ctx context.Context, userName string) (*User, error) {
	userID, err := svc.sf.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID: %w", err)
	}
	user, err := svc.userRepo.CreateUser(ctx, &User{
		ID:   userID,
		Name: userName,
	})
	if err != nil {
		return nil, fmt.Errorf("error create user %d: %w", userID, err)
	}
	return user, nil
}

func (svc *UserServiceImpl) SetUserSession(ctx context.Context, uid uint64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("eror create sid: %w", err)
	}
	sid := base64.URLEncoding.EncodeToString(b)
	if err := svc.userRepo.SetUserSession(ctx, uid, sid); err != nil {
		return "", fmt.Errorf("eror set sid for user %d: %w", uid, err)
	}
	return sid, nil
}

func (svc *UserServiceImpl) GetUserByID(ctx context.Context, uid uint64) (*User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", uid, err)
	}
	return user, nil
}

func (svc *UserServiceImpl) GetUserIDBySession(ctx context.Context, sid string) (uint64, error) {
	userID, err := svc.userRepo.GetUserIDBySession(ctx, sid)
	if err != nil {
		return 0, fmt.Errorf("error get user id by sid %s: %w", sid, err)
	}
	return userID, nil
}

func (svc *UserServiceImpl) GetOrCreateGoogleUser(ctx context.Context, email, userName string) (*User, error) {
	var user *User
	var err error
	user, err = svc.userRepo.GetGoogleUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("error get google user %s: %w", email, err)
		}
		userID, err := svc.sf.NextID()
		if err != nil {
			return nil, fmt.Errorf("error create snowflake ID: %w", err)
		}
		user, err = svc.userRepo.CreateGoogleUser(ctx, email, &User{
			ID:   userID,
			Name: userName,
		})
		if err != nil {
			return nil, fmt.Errorf("error create google user %s: %w", email, err)
		}
		return user, nil
	}
	return user, nil
}
