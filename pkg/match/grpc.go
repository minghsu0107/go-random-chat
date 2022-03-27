package match

import (
	"github.com/minghsu0107/go-random-chat/pkg/config"
	"github.com/minghsu0107/go-random-chat/pkg/transport"
	"google.golang.org/grpc"
)

var (
	ChatConn *ChatClientConn
	UserConn *UserClientConn
)

type ChatClientConn struct {
	Conn *grpc.ClientConn
}

type UserClientConn struct {
	Conn *grpc.ClientConn
}

func NewChatClientConn(config *config.Config) (*ChatClientConn, error) {
	conn, err := transport.InitializeGrpcClient(config.Match.Grpc.Client.Chat.Endpoint)
	if err != nil {
		return nil, err
	}
	ChatConn = &ChatClientConn{
		Conn: conn,
	}
	return ChatConn, nil
}

func NewUserClientConn(config *config.Config) (*UserClientConn, error) {
	conn, err := transport.InitializeGrpcClient(config.Match.Grpc.Client.User.Endpoint)
	if err != nil {
		return nil, err
	}
	UserConn = &UserClientConn{
		Conn: conn,
	}
	return UserConn, nil
}
