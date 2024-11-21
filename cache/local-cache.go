package cache

import (
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/utils"
	"sync"
	"time"
)

const (
	isArray = 1
	isMap   = 0

	defaultCacheSize         = 2                //默认一组cache量
	defaultFlushTimeInterval = 60 * time.Second //默认间隔60秒
)

var (
	bigCacheMap    sync.Map
	startFlushData bool
	lock           sync.Mutex
)

// MapCache 多本地缓存
type MapCache struct {
	// 分段来存储 数据提供并发访问效率
	mapList       []*mapStruct
	listList      *listStruct
	count         int
	nextFlushTime time.Time //由于只有一个公共循环，存储下一次刷新的时间

	dataType          int                                                             //数据类型，默认为0表示map类型，1为数组类型
	CacheType         string                                                          //可以放一类数据，避免不同类型数据覆盖
	MaxLen            int                                                             //总共多长的数据链
	FlushTimeInterval time.Duration                                                   //刷新间隔时间
	FlushCallback     func(mapData *DataEntry, isAllEmpty bool, willDelete bool) bool //自动清理执行的回调方法,willDelete
	// 这次是否会删除，返回bool表示是否要真实删除
}

func newCache(dataType int, cConfig *MapCache) *MapCache {
	if cConfig == nil {
		cConfig = new(MapCache)
	}

	if cConfig.MaxLen <= 0 {
		cConfig.MaxLen = defaultCacheSize
	}
	if cConfig.CacheType == "" {
		cConfig.CacheType = utils.NewUUID() //随机生成一组
	}
	if cConfig.FlushTimeInterval <= 0 {
		cConfig.FlushTimeInterval = defaultFlushTimeInterval
	}

	cConfig.dataType = dataType
	cConfig.nextFlushTime = time.Now().Add(cConfig.FlushTimeInterval)

	oneMultiCache, ok := bigCacheMap.Load(cConfig.CacheType)
	if ok {
		if multiCacheTemp, ok := oneMultiCache.(*MapCache); ok {
			return multiCacheTemp
		}
	}

	if dataType == isArray {
		cConfig.listList = initListCache()
	} else {
		cConfig.mapList = initMapCache(cConfig.MaxLen)
	}

	bigCacheMap.Store(cConfig.CacheType, cConfig)

	//异步来刷缓存的数据,全局可以使用一个，这样可以减少异步的数量
	if !startFlushData {
		lock.Lock()
		defer lock.Unlock()
		if !startFlushData {
			startFlushData = true
			flushAllExpireData()
		}
	}

	return cConfig
}

// NewMapCache 新增一组cache对象
func NewMapCache(cConfig *MapCache) *MapCache {
	if cConfig == nil {
		cConfig = new(MapCache)
	}
	return newCache(isMap, cConfig)
}

// NewListCache 新增一组cache数组对象
func NewListCache(cConfig *MapCache) *MapCache {
	if cConfig == nil {
		cConfig = new(MapCache)
	}
	return newCache(isArray, cConfig)
}

// 全局唯一，减少异步数量
func flushAllExpireData() {
	goroutines.GoAsyncHandler(func(params ...interface{}) {
		//避免有错误，不能循环了
		for {
			bigCacheMap.Range(func(cacheType, mapCache interface{}) bool {
				m, ok := mapCache.(*MapCache)
				if m == nil || !ok {
					return true
				}

				//如果未达到刷新时间，则直接跳过
				nowTime := time.Now()
				if nowTime.Before(m.nextFlushTime) {
					return true
				}

				goroutines.GoSyncHandler(func(params ...interface{}) {
					var runCount = 0

					if m.dataType == isArray {
						runCountTemp := m.listList.flush(m.FlushCallback)
						runCount += runCountTemp
					} else {
						for _, cacheTemp := range m.mapList {
							runCountTemp := cacheTemp.flush(m.FlushCallback)
							runCount += runCountTemp
						}
					}

					if runCount == 0 {
						//表示整个数组都为空了，当次循环没有一个数据了，来执行一次完成以后的事件
						if m.FlushCallback != nil {
							m.FlushCallback(nil, true, false)
						}
					}
				}, nil)

				m.nextFlushTime = nowTime.Add(m.FlushTimeInterval)

				return true
			})
			time.Sleep(time.Second * time.Duration(1)) //1秒以后继续
		}
	}, nil)
}

//func (m *MapCache) flushExpireData() {
//	goroutines.GoAsyncHandlers(func(params ...interface{}) {
//		//避免有错误，不能循环了
//		for {
//			m, ok := params[0].(*MapCache)
//			if m == nil || !ok {
//				break
//			}
//
//			fmt.Println("flushExpireType:", m.CacheType)
//
//			goroutines.GoSynchroHandlers(func(params ...interface{}) {
//				var runCount = 0
//
//				if m.dataType == isArray {
//					runCountTemp := m.listList.flush(m.FlushCallback)
//					runCount += runCountTemp
//				} else {
//					for _, cacheTemp := range m.mapList {
//						runCountTemp := cacheTemp.flush(m.FlushCallback)
//						runCount += runCountTemp
//					}
//				}
//
//				fmt.Println("flushExpireData:", runCount)
//
//				if runCount == 0 {
//					//表示整个数组都为空了，当次循环没有一个数据了，来执行一次完成以后的事件
//					if m.FlushCallback != nil {
//						m.FlushCallback(nil, true, false)
//					}
//				}
//			}, nil)
//			time.Sleep(time.Second * time.Duration(m.FlushTimeInterval))
//		}
//	}, nil, m)
//}

// Range 循环, false 就退出，true就继续
func (m *MapCache) Range(rangeFun func(dataEntry *DataEntry) bool) {
	if m.dataType == isArray {
		m.listList.ranges(rangeFun)
		return
	}
	for _, cacheTemp := range m.mapList {
		ok := cacheTemp.ranges(rangeFun)
		if !ok {
			break
		}
	}
}

// Del 删除
func (m *MapCache) Del(key string) {
	if m.dataType == isArray {
		m.listList.del(key)
		return
	}
	num := utils.Checksum(key, m.MaxLen)
	m.mapList[num].del(key)
}

// 可能有回调的情况
func (m *MapCache) delByCheck(key string) {
	if m.dataType == isArray {
		retData := m.listList.getEntry(key)
		checkDeleteOne(retData, m.FlushCallback)
		m.listList.del(key)
		return
	}
	num := utils.Checksum(key, m.MaxLen)
	retData := m.mapList[num].getEntry(key)
	checkDeleteOne(retData, m.FlushCallback)
	m.mapList[num].del(key)
}

// Get 加载数据
func (m *MapCache) Get(key string) (interface{}, bool) {
	if m.dataType == isArray {
		retData, ok := m.listList.get(key)
		if !ok {
			m.delByCheck(key)
		}
		return retData, ok
	}

	num := utils.Checksum(key, m.MaxLen)
	retData, ok := m.mapList[num].get(key)
	if !ok {
		m.delByCheck(key)
	}
	return retData, ok
}

// Set 存储数据，并更新对应key的expire time
func (m *MapCache) Set(key string, val interface{}, seconds time.Duration) {
	if m.dataType == isArray {
		m.listList.push(m.CacheType, key, val, seconds)
		return
	}

	num := utils.Checksum(key, m.MaxLen)
	m.mapList[num].set(m.CacheType, key, val, seconds)
}

// Count 存储总数
func (m *MapCache) Count() int {
	if m.dataType == isArray {
		return m.listList.count
	}

	count := 0
	for _, cacheTemp := range m.mapList {
		count += cacheTemp.count
	}
	return count
}
