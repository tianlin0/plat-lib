package cache

import "time"

var (
	defaultLocalCache CommCache = NewLocalCache() //本地缓存
)

type commCache struct {
	cCache    CommCache
	isDefault bool //是否是默认的，避免重复提交
}

// New 新建
func New(con ...CommCache) *commCache {
	com := new(commCache)
	if len(con) > 0 {
		com.cCache = con[0]
		return com
	}
	com.isDefault = true
	com.cCache = defaultLocalCache
	return com
}

// Get 从缓存中取得一个值，如果没有redis则从本地缓存
func (co *commCache) Get(key string) (string, error) {
	ret, err := co.cCache.Get(key)
	if err == nil {
		return ret, nil
	}
	if co.isDefault {
		return "", err
	}
	ret2, err2 := defaultLocalCache.Get(key)
	if ret2 == "" || err2 != nil {
		return "", err
	}
	return ret2, nil
}

// Set timeout为秒
func (co *commCache) Set(key, val string, timeout time.Duration) (bool, error) {
	ret, err := co.cCache.Set(key, val, timeout)
	if err == nil {
		return ret, nil
	}
	if co.isDefault {
		return false, err
	}
	ret2, err2 := defaultLocalCache.Set(key, val, timeout)
	if !ret2 || err2 != nil {
		return false, err
	}
	return true, nil
}

// Del 从缓存中删除一个key
func (co *commCache) Del(key string) (bool, error) {
	ret, err := co.cCache.Del(key)
	if err == nil {
		return ret, nil
	}
	if co.isDefault {
		return false, err
	}
	ret2, err2 := defaultLocalCache.Del(key)
	if !ret2 || err2 != nil {
		return false, err
	}
	return true, nil
}
