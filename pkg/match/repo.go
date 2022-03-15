package match

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	chatpb "github.com/minghsu0107/go-random-chat/internal/proto_gen/chat"
	"github.com/minghsu0107/go-random-chat/pkg/infra"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
)

var (
	matchPubSubTopic = "rc:match:pubsub"
	userWaitList     = "rc:userwait"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type ChannelRepo interface {
	CreateChannel(ctx context.Context) (uint64, error)
}

type UserRepo interface {
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
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

func NewChannelRepo(conn *ChatClientConn) ChannelRepo {
	return &ChannelRepoImpl{
		createChannel: transport.NewGrpcEndpoint(
			conn.Conn,
			"chat",
			"chat.ChannelService",
			"CreateChannel",
			&chatpb.CreateChannelResponse{},
		),
	}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context) (uint64, error) {
	res, err := repo.createChannel(ctx, &chatpb.CreateChannelRequest{})
	if err != nil {
		return 0, err
	}
	return res.(*chatpb.CreateChannelResponse).ChannelId, nil
}

type UserRepoImpl struct {
	getUser          endpoint.Endpoint
	addUserToChannel endpoint.Endpoint
}

func NewUserRepo(conn *ChatClientConn) UserRepo {
	return &UserRepoImpl{
		getUser: transport.NewGrpcEndpoint(
			conn.Conn,
			"chat",
			"chat.UserService",
			"GetUser",
			&chatpb.GetUserResponse{},
		),
		addUserToChannel: transport.NewGrpcEndpoint(
			conn.Conn,
			"chat",
			"chat.UserService",
			"AddUserToChannel",
			&chatpb.AddUserResponse{},
		),
	}
}

func (repo *UserRepoImpl) GetUserByID(ctx context.Context, userID uint64) (*User, error) {
	res, err := repo.getUser(ctx, &chatpb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	pbUser := res.(*chatpb.GetUserResponse)
	return &User{
		ID:   pbUser.User.Id,
		Name: pbUser.User.Name,
	}, nil
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
}

func NewMatchingRepo(r infra.RedisCache) MatchingRepo {
	return &MatchingRepoImpl{r}
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
	return repo.r.Publish(ctx, matchPubSubTopic, result.Encode())
}
