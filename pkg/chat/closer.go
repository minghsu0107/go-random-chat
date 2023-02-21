package chat

import (
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

type InfraCloser struct{}

func NewInfraCloser() *InfraCloser {
	return &InfraCloser{}
}

func (closer *InfraCloser) Close() error {
	if err := ForwarderConn.Conn.Close(); err != nil {
		return err
	}
	if err := UserConn.Conn.Close(); err != nil {
		return err
	}
	infra.CassandraSession.Close()
	return infra.RedisClient.Close()
}
