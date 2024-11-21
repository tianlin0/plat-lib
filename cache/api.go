package cache

import (
	"fmt"
	"time"
)

// CommCache 公共缓存接口
type CommCache interface {
	Get(key string) (string, error)
	Set(key, val string, timeout time.Duration) (bool, error)
	Del(key string) (bool, error)
}

// GetNsKey 获取namespace下的key，规范化
func getNsKey(ns string, key string) string {
	if ns != "" {
		return fmt.Sprintf("{%s}%s", ns, key)
	}
	return key
}

// NsGet xxx
func NsGet(co CommCache, ns string, key string) (string, error) {
	return co.Get(getNsKey(ns, key))
}

// NsSet xxx
func NsSet(co CommCache, ns string, key, val string, timeout time.Duration) (bool, error) {
	return co.Set(getNsKey(ns, key), val, timeout)
}

// NsDel xxx
func NsDel(co CommCache, ns string, key string) (bool, error) {
	return co.Del(getNsKey(ns, key))
}
