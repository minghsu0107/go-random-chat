package infra

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/minghsu0107/go-random-chat/pkg/common"
	"github.com/minghsu0107/go-random-chat/pkg/config"
)

var (
	RedisClient redis.UniversalClient
	//ErrRedisUnlockFail is redis unlock fail error
	ErrRedisUnlockFail = errors.New("redis unlock fail")
	// ErrRedisPipelineCmdNotFound is redis command not found error
	ErrRedisPipelineCmdNotFound = errors.New("redis pipeline command not found; supports only SET and DELETE")

	expirationHour int64
	expiration     time.Duration
)

// RedisCache is the interface of redis cache
type RedisCache interface {
	Get(ctx context.Context, key string, dst interface{}) (bool, error)
	Set(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
	HGet(ctx context.Context, key, field string, dst interface{}) (bool, error)
	HMGet(ctx context.Context, key string, fields []string) ([]interface{}, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key, field string) error
	RPush(ctx context.Context, key string, val interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Publish(ctx context.Context, topic string, payload interface{}) error
	ZPopMinOrAddOne(ctx context.Context, key string, score float64, member interface{}) (bool, string, error)
	ZRemOne(ctx context.Context, key string, member interface{}) error
	HGetIfKeyExists(ctx context.Context, key, field string, dst interface{}) (bool, bool, error)
	ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error
	GetMutex(name string) *redsync.Mutex
}

// RedisCacheImpl is the redis cache client type
type RedisCacheImpl struct {
	client redis.UniversalClient
	rs     *redsync.Redsync
}

// RedisOpType is the redis operation type
type RedisOpType int

const (
	// DELETE represents delete operation
	DELETE RedisOpType = iota
	HSETONE
	RPUSH
)

// RedisPayload is a abstract interface for payload type
type RedisPayload interface {
	Payload()
}

// RedisDeletePayload is the payload type of delete method
type RedisDeletePayload struct {
	RedisPayload
	Key string
}

type RedisHsetOnePayload struct {
	RedisPayload
	Key   string
	Field string
	Val   interface{}
}

type RedisRpushPayload struct {
	RedisPayload
	Key string
	Val interface{}
}

// Payload implements abstract interface
func (RedisDeletePayload) Payload()  {}
func (RedisHsetOnePayload) Payload() {}
func (RedisRpushPayload) Payload()   {}

// RedisCmd represents an operation and its payload
type RedisCmd struct {
	OpType  RedisOpType
	Payload RedisPayload
}

// RedisPipelineCmd is redis pipeline command type
type RedisPipelineCmd struct {
	OpType RedisOpType
	Cmd    interface{}
}

func NewRedisClient(config *config.Config) (redis.UniversalClient, error) {
	expirationHour = config.Redis.ExpirationHour
	expiration = time.Duration(expirationHour) * time.Hour
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          common.GetServerAddrs(config.Redis.Addrs),
		Password:       config.Redis.Password,
		ReadOnly:       true,
		RouteByLatency: true,
		MinIdleConns:   config.Redis.MinIdleConn,
		PoolSize:       config.Redis.PoolSize,
		ReadTimeout:    time.Duration(config.Redis.ReadTimeoutMilliSecond) * time.Millisecond,
		WriteTimeout:   time.Duration(config.Redis.WriteTimeoutMilliSecond) * time.Millisecond,
		PoolTimeout:    60 * time.Second,
	})
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	RedisClient.AddHook(redisotel.NewTracingHook())
	return RedisClient, nil
}

// NewRedisCache is the factory of redis cache
func NewRedisCache(client redis.UniversalClient) RedisCache {
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &RedisCacheImpl{
		client: client,
		rs:     rs,
	}
}

// Get returns true if the key already exists and set dst to the corresponding value
func (rc *RedisCacheImpl) Get(ctx context.Context, key string, dst interface{}) (bool, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		json.Unmarshal([]byte(val), dst)
	}
	return true, nil
}

// Set sets a key-value pair
func (rc *RedisCacheImpl) Set(ctx context.Context, key string, val interface{}) error {
	if err := rc.client.Set(ctx, key, val, expiration).Err(); err != nil {
		return err
	}
	return nil
}

// Delete deletes a key
func (rc *RedisCacheImpl) Delete(ctx context.Context, key string) error {
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCacheImpl) HGet(ctx context.Context, key, field string, dst interface{}) (bool, error) {
	val, err := rc.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		json.Unmarshal([]byte(val), dst)
	}
	return true, nil
}

func (rc *RedisCacheImpl) HMGet(ctx context.Context, key string, fields []string) ([]interface{}, error) {
	return rc.client.HMGet(ctx, key, fields...).Result()
}

func (rc *RedisCacheImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rc.client.HGetAll(ctx, key).Result()
}

func (rc *RedisCacheImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rc.client.HSet(ctx, key, values).Err()
}

func (rc *RedisCacheImpl) HDel(ctx context.Context, key, field string) error {
	return rc.client.HDel(ctx, key, field).Err()
}

func (rc *RedisCacheImpl) RPush(ctx context.Context, key string, val interface{}) error {
	return rc.client.RPush(ctx, key, val).Err()
}

func (rc *RedisCacheImpl) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rc.client.LRange(ctx, key, start, stop).Result()
}

func (rc *RedisCacheImpl) Publish(ctx context.Context, topic string, payload interface{}) error {
	return rc.client.Publish(ctx, topic, payload).Err()
}

var zPopMinOrAddOne = redis.NewScript(`
local key = KEYS[1]
local score = ARGV[1]
local member = ARGV[2]
local popmembers = {}

local existed_score = redis.call("ZSCORE", key, member)
if existed_score then
  return ""
end

popmembers = redis.call("ZPOPMIN", key)
if popmembers[1] then
  return popmembers[1]
end

redis.call("ZADD", key, score, member)
return ""
`)

func (rc *RedisCacheImpl) ZPopMinOrAddOne(ctx context.Context, key string, score float64, member interface{}) (bool, string, error) {
	poppedMember, err := zPopMinOrAddOne.Run(ctx, rc.client, []string{key}, score, member).Text()
	if err != nil {
		return false, "", err
	}
	return (poppedMember != ""), poppedMember, nil
}
func (rc *RedisCacheImpl) ZRemOne(ctx context.Context, key string, member interface{}) error {
	return rc.client.ZRem(ctx, key, member).Err()
}

var hgetIfKeyExists = redis.NewScript(`
local key = KEYS[1]
local field = ARGV[1]

if redis.call("EXISTS", key) == 0 then
  return ""
end

return redis.call("HGET", key, field)
`)

func (rc *RedisCacheImpl) HGetIfKeyExists(ctx context.Context, key, field string, dst interface{}) (bool, bool, error) {
	val, err := hgetIfKeyExists.Run(ctx, rc.client, []string{key}, field).Text()
	if err == redis.Nil {
		return true, false, nil
	} else if err != nil {
		return false, false, err
	} else if val == "" {
		return false, false, nil
	} else {
		json.Unmarshal([]byte(val), dst)
	}
	return true, true, nil
}

func (rc *RedisCacheImpl) ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error {
	pipe := rc.client.Pipeline()
	var pipelineCmds []RedisPipelineCmd
	for _, cmd := range *cmds {
		switch cmd.OpType {
		case DELETE:
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: DELETE,
				Cmd:    pipe.Del(ctx, cmd.Payload.(RedisDeletePayload).Key),
			})
		case HSETONE:
			payload := cmd.Payload.(RedisHsetOnePayload)
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: HSETONE,
				Cmd:    pipe.HSet(ctx, payload.Key, payload.Field, payload.Val),
			})
		case RPUSH:
			payload := cmd.Payload.(RedisRpushPayload)
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: RPUSH,
				Cmd:    pipe.RPush(ctx, payload.Key, payload.Val),
			})
		default:
			return ErrRedisPipelineCmdNotFound
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	for _, executedCmd := range pipelineCmds {
		switch executedCmd.OpType {
		case DELETE:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		case HSETONE:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		case RPUSH:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rc *RedisCacheImpl) GetMutex(name string) *redsync.Mutex {
	return rc.rs.NewMutex(name, redsync.WithExpiry(3*time.Second))
}
