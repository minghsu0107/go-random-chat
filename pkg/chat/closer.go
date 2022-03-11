package chat

import "github.com/minghsu0107/go-random-chat/pkg/common"

type InfraCloser struct{}

func NewInfraCloser() common.InfraCloser {
	return &InfraCloser{}
}

func (closer *InfraCloser) Close() error {
	return RedisClient.Close()
}
