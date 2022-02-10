package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
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

func init() {
	var err error
	expirationHour, err = strconv.ParseInt(getenv("REDIS_EXPIRATION_HOURS", "24"), 10, 64)
	if err != nil {
		panic(err)
	}
	expiration = time.Duration(expirationHour) * time.Hour
}

// RedisCache is the interface of redis cache
type RedisCache interface {
	Get(ctx context.Context, key string, dst interface{}) (bool, error)
	Set(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
	HGet(ctx context.Context, key, field string, dst interface{}) (bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key, field string) error
	RPush(ctx context.Context, key string, val interface{}) error
	LLen(ctx context.Context, key string) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Publish(ctx context.Context, topic string, payload interface{}) error
	ZPopMinOrAddOne(ctx context.Context, key string, score float64, member interface{}) (bool, string, error)
	ZRemOne(ctx context.Context, key string, member interface{}) error
	ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error
}

// RedisCacheImpl is the redis cache client type
type RedisCacheImpl struct {
	client redis.UniversalClient
}

// RedisOpType is the redis operation type
type RedisOpType int

const (
	// SET represents set operation
	SET RedisOpType = iota
	// DELETE represents delete operation
	DELETE
)

// RedisPayload is a abstract interface for payload type
type RedisPayload interface {
	Payload()
}

// RedisSetPayload is the payload type of set method
type RedisSetPayload struct {
	RedisPayload
	Key string
	Val interface{}
}

// RedisDeletePayload is the payload type of delete method
type RedisDeletePayload struct {
	RedisPayload
	Key string
}

// Payload implements abstract interface
func (RedisSetPayload) Payload() {}

// Payload implements abstract interface
func (RedisDeletePayload) Payload() {}

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

func NewRedisClient() (redis.UniversalClient, error) {
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         getServerAddrs(getenv("REDIS_ADDRS", "localhost:6379")),
		Password:      getenv("REDIS_PASSWORD", ""),
		ReadOnly:      true,
		RouteRandomly: true,
	})
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	return RedisClient, nil
}

// NewRedisCache is the factory of redis cache
func NewRedisCache(client redis.UniversalClient) RedisCache {
	return &RedisCacheImpl{
		client: client,
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

func (rc *RedisCacheImpl) LLen(ctx context.Context, key string) (int64, error) {
	return rc.client.LLen(ctx, key).Result()
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

func (rc *RedisCacheImpl) ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error {
	pipe := rc.client.Pipeline()
	var pipelineCmds []RedisPipelineCmd
	for _, cmd := range *cmds {
		switch cmd.OpType {
		case SET:
			strVal, err := json.Marshal(cmd.Payload.(RedisSetPayload).Val)
			if err != nil {
				return err
			}
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: SET,
				Cmd:    pipe.Set(ctx, cmd.Payload.(RedisSetPayload).Key, strVal, expiration),
			})
		case DELETE:
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: DELETE,
				Cmd:    pipe.Del(ctx, cmd.Payload.(RedisDeletePayload).Key),
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
		case SET:
			if err := executedCmd.Cmd.(*redis.StatusCmd).Err(); err != nil {
				return err
			}
		case DELETE:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func getServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
