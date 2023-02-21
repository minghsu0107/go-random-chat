package forwarder

import (
	"github.com/minghsu0107/go-random-chat/pkg/infra"
)

type InfraCloser struct{}

func NewInfraCloser() *InfraCloser {
	return &InfraCloser{}
}

func (closer *InfraCloser) Close() error {
	return infra.RedisClient.Close()
}
