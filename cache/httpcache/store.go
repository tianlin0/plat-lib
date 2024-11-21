package httpcache

import (
	"context"
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/internal/gmlock"
	"github.com/tianlin0/plat-lib/internal/gocache/lib/store"
	"time"
)

/*
 * 避免http并发量大时，造成后端数据库访问压力大，而缓慢，进行缓存读取
 * 每次如果命中以后，然后会执行 ExecuteGetDataHandle 更新缓存，这样可以达到实时更新的效果
 */

var (
	//锁
	//lockCacheStore store.StoreInterface
	//oneDataTypeCacheMap        = cmap.New() //保存一组数据，这样可以保证避免CacheType取名相同的，被覆盖了
	storeListCacheMap = cmap.New() //保存store的列表，避免重复创建
	saveDataTypeMap   = cmap.New() //保存每个对象的类型，避免错误初始化
	closed            = false
	//maxMemCacheMb     uint64 = 0 //设置最大使用内存字节数，避免内存爆了
	gmLocker = gmlock.New()
)

//初始化Lock
//func init() {
//	lockCacheStore = gocache_store.NewGoCache(gocache.New(5*time.Minute, 5*time.Minute))
//}

func getStoreCacheKey(namespace string, cacheKey string) string {
	return fmt.Sprintf("{%s}%s", namespace, cacheKey)
}

// 单个获取内容
func getDataFromStore[V any](ctx context.Context, storeList []store.StoreInterface, storeKey string) (value *cacheData[V], err error) {
	var lastErr error
	for _, oneFactory := range storeList {
		one, err := oneFactory.Get(ctx, storeKey)
		if err == nil && one != nil {
			if oneData, ok := one.(*cacheData[V]); ok {
				return oneData, nil
			}
		}
		if err != nil {
			lastErr = err
		}
	}
	return value, lastErr
}

// 根据 store 取得数据
func multiGetData[V any](ctx context.Context, storeList []store.StoreInterface, namespace string, cacheKey string, timeout time.Duration) (value *cacheData[V], err error) {
	if storeList == nil || len(storeList) == 0 {
		return value, fmt.Errorf("multiGetData storeList empty")
	}
	storeKey := getStoreCacheKey(namespace, cacheKey)
	if timeout > 0 {
		value, err = goroutines.RunWithTimeout(timeout, func() (*cacheData[V], error) {
			return getDataFromStore[V](ctx, storeList, storeKey)
		})
		return value, err
	}
	return getDataFromStore[V](ctx, storeList, storeKey)
}

// 根据 store 设置数据
func multiSetData[V any](ctx context.Context, storeList []store.StoreInterface, namespace string, cacheKey string, dataValue V, expiration time.Duration) (bool, error) {
	if storeList == nil || len(storeList) == 0 {
		return false, fmt.Errorf("multiSetData storeList empty")
	}

	if closed {
		return false, fmt.Errorf("closed")
	}

	//_, err := checkNamespaceDataType(namespace, dataValue)
	//if err != nil {
	//	logs.CtxLogger(ctx).Error("multiSetData executeHandle data type error:", err)
	//}

	// TODO: 如果是内存存储，则设置一个内存上限，避免把内存打爆了
	//var memStats runtime.MemStats
	//runtime.ReadMemStats(&memStats)
	//
	//var maxMemCacheByte uint64 = 0
	//if maxMemCacheMb > 0 {
	//	maxMemCacheByte = maxMemCacheMb * 1024 * 1204
	//}
	//if maxMemCacheByte > 0 {
	//	if memStats.Alloc > maxMemCacheByte {
	//		return false, fmt.Errorf("memCache max:%d, %d", memStats.Alloc, maxMemCacheByte)
	//	}
	//} else {
	//	//按内存的占用比例,超过90%就停止
	//	remainingPercentage := (float64(memStats.HeapSys) - float64(memStats.HeapIdle)) / float64(memStats.HeapSys) * 100
	//	logs.CtxLogger(ctx).Warn("memcache:", fmt.Sprintf("%.2f%%", remainingPercentage))
	//}

	storeKey := getStoreCacheKey(namespace, cacheKey)
	var lastErr error
	newCacheData := &cacheData[V]{
		data:           dataValue,
		createTime:     time.Now(),
		expirationTime: time.Now().Add(expiration),
	}
	for _, oneFactory := range storeList {
		err := oneFactory.Set(ctx, storeKey, newCacheData, func(o *store.Options) {
			o.Expiration = expiration
		})
		if err == nil {
			return true, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	return false, lastErr
}

// 根据 store 删除数据
func multiDelData(ctx context.Context, storeList []store.StoreInterface, namespace string, cacheKey string) (bool, error) {
	if storeList == nil || len(storeList) == 0 {
		return false, fmt.Errorf("multiSetData storeList empty")
	}

	storeKey := getStoreCacheKey(namespace, cacheKey)
	var lastErr error
	for _, oneFactory := range storeList {
		err := oneFactory.Delete(ctx, storeKey)
		if err == nil {
			return true, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	return false, lastErr
}

//func checkNamespaceDataType(namespace string, dataValue interface{}) (bool, error) {
//	if value, ok := oneDataTypeCacheMap.Get(namespace); ok {
//		if reflect.TypeOf(dataValue) != value {
//			return false, fmt.Errorf("checkNamespaceDataType data type error:%s, %v", namespace, dataValue)
//		}
//		return true, nil
//	}
//	if dataValue != nil {
//		oneDataTypeCacheMap.Set(namespace, reflect.TypeOf(dataValue)) //存储类型
//	}
//	return false, nil
//}

// 同一个key的访问次数
func getLockCacheKey(namespace string, cacheKey string) string {
	return fmt.Sprintf("{%s}{lock-key}%s", namespace, cacheKey)
}

//// 为了避免同一个请求同时执行查询多次，需要对当前的key创建通道，保证数量
//func getLockMutex(ctx context.Context, namespace string, cacheKey string) (currentLockerKey string, currentLockerMutex sync.Mutex, isNew bool, err error) {
//	currentLockerKey = getLockCacheKey(namespace, cacheKey, 0)
//
//	lockData, err := lockCacheStore.Get(ctx, currentLockerKey)
//	if err == nil && lockData != nil {
//		if lockMutex, ok := lockData.(sync.Mutex); ok {
//			return currentLockerKey, lockMutex, false, nil
//		}
//	}
//	currentLockerMutex = sync.Mutex{}
//	err = lockCacheStore.Set(ctx, currentLockerKey, currentLockerMutex)
//	return currentLockerKey, currentLockerMutex, true, err
//}
//
//func delLockMutex(ctx context.Context, currentLockerKey string) error {
//	return lockCacheStore.Delete(ctx, currentLockerKey)
//}
