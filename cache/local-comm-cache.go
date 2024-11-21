package cache

import (
	"github.com/tianlin0/plat-lib/conv"
	"time"
)

var (
	localCache = NewMapCache(&MapCache{
		MaxLen: 5,
	})
)

type localCommCache struct {
	r *MapCache
}

// NewLocalCache 新建
func NewLocalCache(con ...*MapCache) *localCommCache {
	com := new(localCommCache)
	if len(con) > 0 {
		com.r = NewMapCache(con[0])
		return com
	}
	com.r = localCache
	return com
}

// Get 从缓存中取得一个值，如果没有redis则从本地缓存
func (co *localCommCache) Get(key string) (string, error) {
	retData, ok := co.r.Get(key)
	if ok {
		retDataStr := conv.String(retData)
		if retDataStr != "" {
			return retDataStr, nil
		}
	}
	return "", nil
}

// Set timeout为秒
func (co *localCommCache) Set(key, val string, timeout time.Duration) (bool, error) {
	co.r.Set(key, val, timeout)
	return true, nil
}

// Del 从缓存中删除一个key
func (co *localCommCache) Del(key string) (bool, error) {
	co.r.Del(key)
	return true, nil
}
