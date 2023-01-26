package common

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const rateLimitRedisKeyPrefix = "rc:ratelimit"

type RateLimiter struct {
	rc         redis.UniversalClient
	rate       int
	burst      int
	expiration time.Duration
}

var rateLimitScript = redis.NewScript(`
local tokens_key = KEYS[1]
local timestamp_key = KEYS[2]
local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local ttl = math.floor(tonumber(ARGV[5]))
local last_tokens = tonumber(redis.call("get", tokens_key))
if last_tokens == nil then
    last_tokens = capacity
end
local last_refreshed = tonumber(redis.call("get", timestamp_key))
if last_refreshed == nil then
    last_refreshed = 0
end
local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end
redis.call("setex", tokens_key, ttl, new_tokens)
redis.call("setex", timestamp_key, ttl, now)
return { allowed, new_tokens }
`)

// NewRateLimiter returns a new Limiter that allows events up to rate r
// and permits bursts of at most b tokens
func NewRateLimiter(rc redis.UniversalClient, rate, burst int, expiration time.Duration) *RateLimiter {
	return &RateLimiter{
		rc:         rc,
		rate:       rate,
		burst:      burst,
		expiration: expiration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	return rl.AllowN(ctx, key, time.Now(), 1)
}

func (rl *RateLimiter) AllowN(ctx context.Context, key string, now time.Time, n int) (bool, error) {
	reservation, err := rl.reserveN(ctx, Join(rateLimitRedisKeyPrefix, ":", key), now, n)
	if err != nil {
		return false, err
	}
	return reservation.ok, nil
}

type Reservation struct {
	ok     bool
	tokens int
}

func (rl *RateLimiter) reserveN(ctx context.Context, key string, now time.Time, n int) (*Reservation, error) {
	// force key to be hashed to the same slot
	tokenKey := Join("{", key, "}", ":tokens")
	timestampKey := Join("{", key, "}", ":ts")
	result, err := rateLimitScript.Run(ctx, rl.rc, []string{tokenKey, timestampKey}, float64(rl.rate), rl.burst, now.Unix(), n, rl.expiration.Seconds()).Result()
	if err != nil {
		return nil, err
	}

	rs, ok := result.([]interface{})
	if !ok {
		return nil, err
	}
	newTokens, _ := rs[1].(int64)
	return &Reservation{
		ok:     rs[0] == int64(1),
		tokens: int(newTokens),
	}, nil
}
