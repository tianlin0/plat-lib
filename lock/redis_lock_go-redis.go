package lock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	DefaultExpireTime = 10 * time.Second // 默认过期时间10s
	DefaultKeyFront   = "{redis-lock}"
)

// redisLock
type redisLock struct {
	redisClient *redis.Client
	key         string
	value       string
	expire      time.Duration
}

func getLockerKeyName(key string) string {
	return fmt.Sprintf("%s%s", DefaultKeyFront, key)
}

// NewRedisLock 新的锁
func NewRedisLock(redisClient *redis.Client, key string, expire time.Duration) (*redisLock, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	v := base64.StdEncoding.EncodeToString(b)

	//redis时间不能太短，避免大量的redis操作
	if expire < DefaultExpireTime {
		expire = DefaultExpireTime
	}

	return &redisLock{
		redisClient: redisClient,
		key:         getLockerKeyName(key),
		value:       v,
		expire:      expire,
	}, nil
}

// Lock 上锁
func (l *redisLock) Lock(ctx context.Context) (bool, error) {
	return l.redisClient.SetNX(ctx, l.key, l.value, l.expire).Result()
}

// UnLock 解锁
func (l *redisLock) UnLock(ctx context.Context) (bool, error) {
	script := "if redis.call(\"GET\", KEYS[1]) == ARGV[1] then return redis.call(\"DEL\", KEYS[1]) else return 0 end"
	return l.redisClient.Eval(ctx, script, []string{"1"}, l.key, l.value).Bool()
}
