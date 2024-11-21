package cache

import (
	"github.com/tianlin0/plat-lib/conv"
	"sync"
	"time"
)

type mapStruct struct {
	m sync.Map

	count int //总数量
	lock  sync.Mutex
}

// DataEntry 使用，如果过期时间本来就小于创建时间，则表示永远不过期
type DataEntry struct {
	CTime  time.Time //创建时间
	Expire time.Time //过期时间
	Type   string    //同一类名字
	Key    string
	Value  interface{}
}

func initMapCache(maxLen int) []*mapStruct {
	cacheList := make([]*mapStruct, maxLen)
	for i := 0; i < len(cacheList); i++ {
		cacheList[i] = new(mapStruct) //初始化
	}
	return cacheList
}

// Set 设置
func (c *mapStruct) set(cacheType string, key string, value interface{}, seconds time.Duration) {
	//如果seconds < 0 表示永远不过期
	entry := &DataEntry{
		Expire: time.Now().Add(seconds),
		CTime:  time.Now(),
		Type:   cacheType,
		Key:    key,
		Value:  value,
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.m.Load(key)
	if !ok { //覆盖的话，则不执行
		c.count = c.count + 1
	}
	c.m.Store(key, entry)
}

// Get 获取, bool 表示可用
func (c *mapStruct) get(key string) (retData interface{}, canUse bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}
	e, ok := v.(*DataEntry)
	if !ok {
		c.del(key)
		return nil, false
	}
	return e.Value, dataValidate(e)
}

func (c *mapStruct) getEntry(key string) *DataEntry {
	v, ok := c.m.Load(key)
	if !ok {
		return nil
	}
	e, ok := v.(*DataEntry)
	if !ok {
		return nil
	}
	return e
}

// Del 删除
func (c *mapStruct) del(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.m.Load(key)
	if ok {
		c.count = c.count - 1
	}
	c.m.Delete(key)
	return ok
}

// ranges
// rangeBack 返回false 跳出循环
func (c *mapStruct) ranges(rangeBack func(dataEntry *DataEntry) bool) bool {
	retBool := true
	var deleteKeyList []interface{}
	c.m.Range(func(keyOld, valueOld interface{}) bool {
		e, ok := valueOld.(*DataEntry)
		if !ok {
			if deleteKeyList == nil {
				deleteKeyList = make([]interface{}, 0)
			}
			deleteKeyList = append(deleteKeyList, keyOld)
			return true
		}
		retBool = rangeBack(e)
		return retBool
	})
	if deleteKeyList != nil {
		// 在循环外删，避免死锁的情况
		for _, keyTemp := range deleteKeyList {
			c.m.Delete(keyTemp)
		}
	}
	return retBool
}

// 清除过期的，最大的数量maxLen，多于这个数的，则会清空时间最久的，返回目前执行后的数据数量
// deleteCallback 返回true 表示删除
func (c *mapStruct) flush(deleteCallback func(dataEntry *DataEntry, isEmpty bool, isDelete bool) bool) int {
	var runListCount = 0
	c.m.Range(func(key, value interface{}) bool {
		runListCount++
		e, ok := value.(*DataEntry)
		keyStr := conv.String(key)
		if !ok {
			e = nil
		}
		deleteKey := checkDeleteOne(e, deleteCallback)
		//过期将要删除，然后手动函数也返回能删除，则该项目会被删除掉
		if deleteKey {
			c.del(keyStr)
			runListCount--
		}
		return true
	})
	return runListCount
}

func checkDeleteOne(e *DataEntry, deleteCallback func(dataEntry *DataEntry, isEmpty bool, isDelete bool) bool) bool {
	if e == nil {
		return true
	}

	isWillDelete := false
	okData := dataValidate(e)
	if !okData {
		isWillDelete = true
	}

	userDelete := false //用户返回是否要删除
	canDelete := true
	if deleteCallback != nil {
		canDelete = deleteCallback(e, false, isWillDelete)
		userDelete = true
	}
	deleteKey := false
	if userDelete { //如果有用户接口返回要求删除的话
		if canDelete {
			deleteKey = true
		}
	} else {
		if isWillDelete {
			deleteKey = true
		}
	}
	return deleteKey
}
