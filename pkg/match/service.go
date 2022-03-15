package match

import (
	"context"

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
	return svc.userRepo.GetUserByID(ctx, uid)
}
func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {
	return svc.userRepo.AddUserToChannel(ctx, channelID, userID)
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
		return nil, err
	}
	if matched {
		newChannelID, err := svc.chanRepo.CreateChannel(ctx)
		if err != nil {
			return nil, err
		}
		accessToken, err := common.NewJWT(newChannelID)
		if err != nil {
			return nil, err
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
	return svc.matchRepo.PublishMatchResult(ctx, result)
}
func (svc *MatchingServiceImpl) RemoveUserFromWaitList(ctx context.Context, userID uint64) error {
	return svc.matchRepo.RemoveFromWaitList(ctx, userID)
}
