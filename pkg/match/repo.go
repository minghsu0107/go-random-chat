package match

import (
	"context"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	chatpb "github.com/minghsu0107/go-random-chat/proto/chat"
	userpb "github.com/minghsu0107/go-random-chat/proto/user"
)

var (
	matchPubSubTopic = "rc.match"
	userWaitList     = "rc:userwait"
)

type ChannelRepo interface {
	CreateChannel(ctx context.Context) (uint64, string, error)
}

type UserRepo interface {
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
	GetUserIDBySession(ctx context.Context, sid string) (uint64, error)
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
}

type MatchingRepo interface {
	PopOrPushWaitList(ctx context.Context, userID uint64) (bool, uint64, error)
	PublishMatchResult(ctx context.Context, result *MatchResult) error
	RemoveFromWaitList(ctx context.Context, userID uint64) error
}

type ChannelRepoImpl struct {
	createChannel endpoint.Endpoint
}

func NewChannelRepoImpl(chatConn *ChatClientConn) *ChannelRepoImpl {
	return &ChannelRepoImpl{
		createChannel: transport.NewGrpcEndpoint(
			chatConn.Conn,
			"chat",
			"chat.ChannelService",
			"CreateChannel",
			&chatpb.CreateChannelResponse{},
		),
	}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context) (uint64, string, error) {
	res, err := repo.createChannel(ctx, &chatpb.CreateChannelRequest{})
	if err != nil {
		return 0, "", err
	}
	return res.(*chatpb.CreateChannelResponse).ChannelId, res.(*chatpb.CreateChannelResponse).AccessToken, nil
}

type UserRepoImpl struct {
	getUserByID        endpoint.Endpoint
	getUserIDBySession endpoint.Endpoint
	addUserToChannel   endpoint.Endpoint
}

func NewUserRepoImpl(userConn *UserClientConn, chatConn *ChatClientConn) *UserRepoImpl {
	return &UserRepoImpl{
		getUserByID: transport.NewGrpcEndpoint(
			userConn.Conn,
			"user",
			"user.UserService",
			"GetUser",
			&userpb.GetUserResponse{},
		),
		getUserIDBySession: transport.NewGrpcEndpoint(
			userConn.Conn,
			"user",
			"user.UserService",
			"GetUserIdBySession",
			&userpb.GetUserIdBySessionResponse{},
		),
		addUserToChannel: transport.NewGrpcEndpoint(
			chatConn.Conn,
			"chat",
			"chat.UserService",
			"AddUserToChannel",
			&chatpb.AddUserResponse{},
		),
	}
}

func (repo *UserRepoImpl) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	res, err := repo.getUserByID(ctx, &userpb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	pbUser := res.(*userpb.GetUserResponse)
	if !pbUser.Exist {
		return nil, ErrUserNotFound
	}
	return &User{
		ID:   pbUser.User.Id,
		Name: pbUser.User.Name,
	}, nil
}

func (repo *UserRepoImpl) GetUserIDBySession(ctx context.Context, sid string) (uint64, error) {
	res, err := repo.getUserIDBySession(ctx, &userpb.GetUserIdBySessionRequest{
		Sid: sid,
	})
	if err != nil {
		return 0, err
	}
	pbUserID := res.(*userpb.GetUserIdBySessionResponse)
	return pbUserID.UserId, nil
}

func (repo *UserRepoImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	_, err := repo.addUserToChannel(ctx, &chatpb.AddUserRequest{
		ChannelId: channelID,
		UserId:    userID,
	})
	if err != nil {
		return err
	}
	return nil
}

type MatchingRepoImpl struct {
	r infra.RedisCache
	p message.Publisher
}

func NewMatchingRepoImpl(r infra.RedisCache, p message.Publisher) *MatchingRepoImpl {
	return &MatchingRepoImpl{r, p}
}
func (repo *MatchingRepoImpl) PopOrPushWaitList(ctx context.Context, userID uint64) (bool, uint64, error) {
	match, peerIDStr, err := repo.r.ZPopMinOrAddOne(ctx, userWaitList, float64(time.Now().Unix()), userID)
	if err != nil {
		return false, 0, err
	}
	if !match {
		return false, 0, nil
	}
	peerID, err := strconv.ParseUint(peerIDStr, 10, 64)
	if err != nil {
		return false, 0, err
	}
	return true, peerID, nil
}
func (repo *MatchingRepoImpl) RemoveFromWaitList(ctx context.Context, userID uint64) error {
	return repo.r.ZRemOne(ctx, userWaitList, userID)
}
func (repo *MatchingRepoImpl) PublishMatchResult(ctx context.Context, result *MatchResult) error {
	return repo.p.Publish(matchPubSubTopic, message.NewMessage(
		watermill.NewUUID(),
		result.Encode(),
	))
}
