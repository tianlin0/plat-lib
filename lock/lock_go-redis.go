package lock

import (
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/cache"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/internal/gmlock"
	"time"
)

// lockRedis 分布式锁 redis-server
func lockRedis(ctx context.Context, key string, timeout time.Duration) (bool, *redisLock, error) {
	redisClient, _ := cache.GetRedisClient(nil)
	if redisClient == nil {
		return false, nil, fmt.Errorf("lock has no redis conn")
	}

	redisLocks, err := NewRedisLock(redisClient, key, timeout)
	if err != nil {
		return false, nil, err
	}
	if redisLocks == nil {
		return false, nil, fmt.Errorf("lockClient is nil")
	}
	retBool, err := redisLocks.Lock(ctx)
	return retBool, redisLocks, err
}

var gmLocker = gmlock.New()

// Lock 加锁
func Lock(key string, callFunction func(), timeout ...time.Duration) bool {
	timeoutExp := DefaultExpireTime
	if timeout != nil || len(timeout) > 0 {
		timeoutExp = timeout[0]
	}

	ctx := context.Background()
	//先用redis
	retBool, client, err := lockRedis(ctx, key, timeoutExp)
	if err == nil && client != nil {
		if !retBool { //表示锁住了
			return true
		}
		defer func(client *redisLock) {
			_, _ = client.UnLock(ctx)
		}(client)

		goroutines.GoSyncHandler(func(params ...interface{}) {
			callFunction()
		}, nil)
		return false
	}

	//内存锁
	lockKey := getLockerKeyName(key)

	gmLocker.Lock(lockKey)
	defer func() {
		gmLocker.Unlock(lockKey)
		gmLocker.Remove(lockKey)
	}()

	goroutines.GoSyncHandler(func(params ...interface{}) {
		callFunction()
	}, nil)

	return false
}
