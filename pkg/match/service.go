package match

import (
	"context"
	"fmt"

	"github.com/minghsu0107/go-random-chat/pkg/common"
)

type UserService interface {
	GetUser(ctx context.Context, uid uint64) (*User, error)
	AddUserToChannel(ctx context.Context, channelID, userID uint64) error
}

type MatchingService interface {
	Match(ctx context.Context, userID uint64) (*MatchResult, error)
	BroadcastMatchResult(ctx context.Context, result *MatchResult) error
	RemoveUserFromWaitList(ctx context.Context, userID uint64) error
}

type UserServiceImpl struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) UserService {
	return &UserServiceImpl{userRepo}
}

func (svc *UserServiceImpl) GetUser(ctx context.Context, uid uint64) (*User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", uid, err)
	}
	return user, nil
}
func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.AddUserToChannel(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error add user %d to channel %d: %w", userID, channelID, err)
	}
	return nil
}

type MatchingServiceImpl struct {
	matchRepo MatchingRepo
	chanRepo  ChannelRepo
}

func NewMatchingService(matchRepo MatchingRepo, chanRepo ChannelRepo) MatchingService {
	return &MatchingServiceImpl{matchRepo, chanRepo}
}
func (svc *MatchingServiceImpl) Match(ctx context.Context, userID uint64) (*MatchResult, error) {
	matched, peerID, err := svc.matchRepo.PopOrPushWaitList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error match user %d: %w", userID, err)
	}
	if matched {
		newChannelID, err := svc.chanRepo.CreateChannel(ctx)
		if err != nil {
			return nil, fmt.Errorf("error create channel: %w", err)
		}
		accessToken, err := common.NewJWT(newChannelID)
		if err != nil {
			return nil, fmt.Errorf("error create JWT: %w", err)
		}
		return &MatchResult{
			Matched:     true,
			UserID:      userID,
			PeerID:      peerID,
			ChannelID:   newChannelID,
			AccessToken: accessToken,
		}, nil
	}
	return &MatchResult{
		Matched: false,
	}, nil
}
func (svc *MatchingServiceImpl) BroadcastMatchResult(ctx context.Context, result *MatchResult) error {
	if err := svc.matchRepo.PublishMatchResult(ctx, result); err != nil {
		return fmt.Errorf("error broadcast match result: %w", err)
	}
	return nil
}
func (svc *MatchingServiceImpl) RemoveUserFromWaitList(ctx context.Context, userID uint64) error {
	if err := svc.matchRepo.RemoveFromWaitList(ctx, userID); err != nil {
		return fmt.Errorf("error remove user %d from wait list: %w", userID, err)
	}
	return nil
}
