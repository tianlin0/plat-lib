package httpcache

import (
	"context"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	gocache "github.com/patrickmn/go-cache"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/internal/gocache/lib/store"
	gocache_store "github.com/tianlin0/plat-lib/internal/gocache/store/go_cache"
	"github.com/tianlin0/plat-lib/internal/gocache/store/lrucache"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"runtime"
	"time"
)

var (
	maxExpiration               = 24 * 7 * time.Hour
	defaultExpiration           = 24 * time.Hour
	defaultCleanupTimes         = 5
	minCleanupInterval          = 10 * time.Minute
	defaultAsyncExecuteDuration = 5 * time.Minute
)

/*
 * 避免http并发量大时，造成后端数据库访问压力大，而缓慢，进行缓存读取
 * 每次如果命中以后，然后会执行 ExecuteGetDataHandle 更新缓存，这样可以达到实时更新的效果
 */
func checkParam[P any, V any](cfg *Config[P, V], fileName string, fileLine int) (*Config[P, V], error) {
	if cfg.Namespace == "" {
		fileList := utils.GetRuntimeCallers(fileName, fileLine, 0, 1)
		if len(fileList) == 0 {
			return nil, fmt.Errorf("NameSpace is empty")
		}
		nameSpace := fmt.Sprintf("%s:%d/%s", fileList[0].FileName, fileList[0].Line, fileList[0].FuncName)
		cfg.Namespace = utils.GetUUID(nameSpace)
		logs.DefaultLogger().Error("Namespace auto create:", cfg.Namespace, nameSpace)
	}

	//Timeout 为0的话，则默认的defaultExpiration=-1，表示不限制时间
	if cfg.Expiration < 0 {
		cfg.Expiration = maxExpiration //最长缓存7天，避免永久缓存
	} else if cfg.Expiration == 0 {
		cfg.Expiration = defaultExpiration //表示未设置，默认1天
	}

	if cfg.CleanupInterval == 0 {
		cfg.CleanupInterval = cfg.Expiration / time.Duration(defaultCleanupTimes) //一个周期内执行清理5次
	}
	//清理不能太频繁，避免运行次数太多
	if cfg.CleanupInterval < minCleanupInterval {
		cfg.CleanupInterval = minCleanupInterval
	}

	if cfg.AsyncExecuteDuration == 0 {
		cfg.AsyncExecuteDuration = defaultAsyncExecuteDuration //默认5分钟之内不进行自动更新
	}

	var err error
	if cfg.GetDataExecuteHandle == nil {
		err = fmt.Errorf("ExecuteGetDataHandle null")
	}
	if cfg.StoreList == nil || len(cfg.StoreList) == 0 {
		if storeList, ok := storeListCacheMap.Get(cfg.Namespace); ok {
			if storeListTemp, ok := storeList.([]store.StoreInterface); ok {
				cfg.StoreList = storeListTemp
			}
		}

		if cfg.StoreList == nil || len(cfg.StoreList) == 0 {
			cfg.StoreList = []store.StoreInterface{
				newDefaultStore(cfg),
			}
		}
	}

	//默认设置第一个
	defaultSet := false
	if storeList, ok := storeListCacheMap.Get(cfg.Namespace); ok {
		if _, ok := storeList.([]store.StoreInterface); ok {
			defaultSet = true
		}
	}
	if !defaultSet {
		storeListCacheMap.Set(cfg.Namespace, cfg.StoreList)
	}

	return cfg, err
}

// 根据参数初始化默认Store
func newDefaultStore[P any, V any](cfg *Config[P, V]) store.StoreInterface {
	if cfg.MaxSize == 0 {
		//不需要设置总数
		//默认用go_cache
		goCacheClient := gocache.New(cfg.Expiration, cfg.CleanupInterval)
		goCacheStore := gocache_store.NewGoCache(goCacheClient)

		return goCacheStore
	}

	lruCacheClient := expirable.NewLRU[string, any](cfg.MaxSize, nil, cfg.Expiration)
	lruCacheStore := lrucache.NewLRUCache[string, any](lruCacheClient)
	return lruCacheStore
}

func New[P any, V any](cfg *Config[P, V]) (HttpCache[P, V], error) {
	fileName := ""
	fileLine := 0
	if cfg.Namespace == "" {
		_, fileName, fileLine, _ = runtime.Caller(1)
	}
	newCfg, err := checkParam(cfg, fileName, fileLine)
	if err != nil {
		logs.DefaultLogger().Error("httpCache New Error:", cfg.Namespace, err)
		if newCfg == nil {
			panic("httpCache New error:" + cfg.Namespace)
			return nil, err
		}
	}

	n := new(cacheIns[P, V])
	n.cfg = newCfg

	return n, err
}

// Close 关闭httpCache
func Close() {
	closed = true
}

func (c *cacheIns[P, V]) needAsyncGetData(ctx context.Context, tempData *cacheData[V]) bool {
	//如果小于0，表示不用实时更新，固定数据，不会变更
	if c.cfg.AsyncExecuteDuration < 0 {
		return false
	}

	isUpdate := false
	duration := time.Now().Sub(tempData.createTime)
	if duration > c.cfg.AsyncExecuteDuration {
		isUpdate = true //如果超时了，则异步更新
	} else {
		//根据程序判断是否需要异步自动更新
		if c.cfg.NeedAsyncExecuteHandle != nil {
			isUpdate = c.cfg.NeedAsyncExecuteHandle(ctx, tempData.data)
		}
	}

	return isUpdate
}

// GetMap 获取多个对象
func (c *cacheIns[P, V]) getMap(ctx context.Context, cacheMapKeys map[string]P) (map[string]V, error) {
	retMap := make(map[string]V)

	if len(cacheMapKeys) == 0 {
		return retMap, fmt.Errorf("cacheKey is empty")
	}

	found := false
	unCacheKey := make(map[string]P, 0)
	asyncCacheKey := make(map[string]P, 0)
	for oneCacheKey, oneCacheParam := range cacheMapKeys {
		if oneCacheKey == "" {
			continue
		}
		found = true
		tempData, err := c.getOneFromCache(ctx, oneCacheKey)
		if err == nil && tempData != nil {
			isUpdate := c.needAsyncGetData(ctx, tempData)
			if isUpdate {
				asyncCacheKey[oneCacheKey] = oneCacheParam
			}
			retMap[oneCacheKey] = tempData.data
			continue
		}
		unCacheKey[oneCacheKey] = oneCacheParam
	}

	//传入了多个空字符串
	if !found {
		return retMap, fmt.Errorf("cacheKeys is empty")
	}

	//判断是否有自动获取数据的接口，没有则直接返回，提高执行效率
	if c.cfg.GetDataExecuteHandle == nil {
		return retMap, nil
	}

	//表示有多个需要重新覆盖
	if len(asyncCacheKey) > 0 {
		goroutines.GoSyncHandler(func(params ...interface{}) {
			//异步有可能指针的原始值已被修改了
			if asyncCacheKeyList, ok := params[0].(map[string]P); ok {
				_, _ = c.lockExecuteListHandle(ctx, false, asyncCacheKeyList)
			}
		}, nil, asyncCacheKey)
	}

	if len(unCacheKey) == 0 {
		return retMap, nil
	}

	newMap, err := c.lockExecuteListHandle(ctx, true, unCacheKey)

	logger := logs.CtxLogger(ctx)

	logger.Info("httpCache lockExecuteListHandle End", err)
	if newMap != nil {
		for k, v := range newMap {
			retMap[k] = v
		}
	}
	for k, v := range retMap {
		if cond.IsNil(v) {
			delete(retMap, k)
		}
	}

	return retMap, err
}

// Get 获取一个对象
func (c *cacheIns[P, V]) Get(ctx context.Context, cacheKey string, getDataParam P) (value V, err error) {
	retMap, err := c.getMap(ctx, map[string]P{
		cacheKey: getDataParam,
	})
	logger := logs.CtxLogger(ctx)
	if err != nil {
		logger.Error("httpCache getMap Error:", err)
		return value, err
	}
	if data, ok := retMap[cacheKey]; ok {
		logger.Debug("httpCache Get Result")
		return data, nil
	}
	logger.Error("httpCache getMap No data")
	return value, nil
}

// Set 外部手动进行设置
func (c *cacheIns[P, V]) Set(ctx context.Context, cacheKey string, dataValue V) bool {
	if cond.IsNil(dataValue) || cacheKey == "" {
		return false
	}
	ret, _ := multiSetData(ctx, c.cfg.StoreList, c.cfg.Namespace, cacheKey, dataValue, c.cfg.Expiration)
	return ret
}
func (c *cacheIns[P, V]) Del(ctx context.Context, cacheKey string) bool {
	if cacheKey == "" {
		return false
	}
	ret, _ := multiDelData(ctx, c.cfg.StoreList, c.cfg.Namespace, cacheKey)
	return ret
}

// Get 获取一个对象
func (c *cacheIns[P, V]) getOneFromCache(ctx context.Context, oneCacheKey string) (cData *cacheData[V], err error) {
	tempData, err := multiGetData[V](ctx, c.cfg.StoreList, c.cfg.Namespace, oneCacheKey, c.cfg.Timeout)
	if err == nil && tempData != nil {
		return tempData, nil
	}
	if err == nil {
		err = fmt.Errorf("no found data")
	}
	return nil, err
}

// lockExecuteListHandle 获取缓存数据的方法
func (c *cacheIns[P, V]) lockExecuteListHandle(ctx context.Context, wait bool, cacheKeyMap map[string]P) (map[string]V, error) {
	retMap := make(map[string]V)
	var getErr error

	for oneCacheKey, oneDataParam := range cacheKeyMap {
		retData, err := c.exeOneFunction(ctx, wait, oneCacheKey, oneDataParam)
		if err != nil {
			getErr = err
			continue
		}
		retMap[oneCacheKey] = retData
	}

	return retMap, getErr
}

func (c *cacheIns[P, V]) exeOneFunction(ctx context.Context, wait bool, oneCacheKey string, getDataParam P) (value V, err error) {
	var oldCacheDataTime time.Time
	{ //取出现在缓存中的创建时间
		tempData, errTemp := c.getOneFromCache(ctx, oneCacheKey)
		if errTemp == nil && tempData != nil {
			oldCacheDataTime = tempData.createTime
		}
	}

	logger := logs.CtxLogger(ctx)

	logger.Info("httpCache exeOneFunction lock Start", oneCacheKey, wait)

	retData, err := c.lock(ctx, oneCacheKey, wait, func(ctx context.Context) (V, error) {
		loggerIn := logs.CtxLogger(ctx)

		loggerIn.Info("httpCache exeOneFunction lock func", oldCacheDataTime, oneCacheKey, getDataParam)

		//这里需要进行二次查询，避免第一次查询成功以后，大量重复执行，只负责从缓存中获取
		if !cond.IsTimeEmpty(oldCacheDataTime) {
			tempData, errTemp := c.getOneFromCache(ctx, oneCacheKey)
			if errTemp == nil && tempData != nil {
				//这里需要对创建时间进行判断，如果时间变更的话，则直接返回
				if tempData.createTime.Sub(oldCacheDataTime) > 0 {
					return tempData.data, nil
				}
			}
		}

		return c.executeHandle(ctx, oneCacheKey, getDataParam)
	})

	logger.Debug("httpCache exeOneFunction lock End", oneCacheKey)

	if err != nil {
		logger.Error("exeOneCache:", err)
		return value, err
	}
	if !cond.IsNil(retData) {
		return retData, nil
	}

	return value, fmt.Errorf("no data:%s", oneCacheKey)
}

func (c *cacheIns[P, V]) executeHandle(ctx context.Context, cacheKey string, getDataParam P) (value V, err error) {
	if c.cfg.GetDataExecuteHandle == nil {
		return value, nil
	}

	logger := logs.CtxLogger(ctx)

	goroutines.GoSyncHandler(func(params ...interface{}) {
		ctx1, _ := params[0].(context.Context)
		cacheKey1, ok2 := params[1].(string)
		getDataParam1, ok3 := params[2].(P)
		if ok2 && ok3 {
			loggerIn := logs.CtxLogger(ctx1)
			loggerIn.Info("executeHandle ExecuteGetDataHandle start:", cacheKey1, getDataParam1)
			value, err = c.cfg.GetDataExecuteHandle(ctx1, cacheKey1, getDataParam1)
		}
	}, nil, ctx, cacheKey, getDataParam)

	//返回为nil，表示不自动设置，可能需要进行外部设置
	if cond.IsNil(value) {
		logger.Error("executeHandle ExecuteGetDataHandle nil:", cacheKey)
		return value, err
	}

	if err == nil {
		c.Set(ctx, cacheKey, value)
	}
	return value, err
}

func (c *cacheIns[P, V]) lock(ctx context.Context, cacheKey string, wait bool, fun func(ctx context.Context) (V, error)) (value V, err error) {
	lockerKey := getLockCacheKey(c.cfg.Namespace, cacheKey)

	if gmLocker.Locked(lockerKey) {
		if !wait {
			return value, fmt.Errorf("currentLockerKey wait: %s, %v", lockerKey, wait)
		}
	}

	logger := logs.CtxLogger(ctx)

	gmLocker.Lock(lockerKey)
	logger.Info("httpCache lock enter1", lockerKey, cacheKey)
	defer func() {
		gmLocker.Unlock(lockerKey)
	}()

	return fun(ctx)
}
