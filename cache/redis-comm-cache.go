package cache

import (
	"context"
	"github.com/tianlin0/plat-lib/conn"
	"time"
)

var (
	ctx = context.Background()
)

type redisCommCache struct {
	r *redisClient
}

// NewRedisCache 新建
func NewRedisCache(con ...*conn.Connect) *redisCommCache {
	var tempCon *conn.Connect
	if len(con) > 0 {
		tempCon = con[0]
	}
	com := new(redisCommCache)
	com.r = NewRedis(tempCon)
	return com
}

// Get 从缓存中取得一个值，如果没有redis则从本地缓存
func (co *redisCommCache) Get(key string) (string, error) {
	return co.r.Get(ctx, key)
}

// Set timeout为秒
func (co *redisCommCache) Set(key, val string, timeout time.Duration) (bool, error) {
	return co.r.Set(ctx, key, val, timeout)
}

// Del 从缓存中删除一个key
func (co *redisCommCache) Del(key string) (bool, error) {
	return co.r.Del(ctx, key)
}
