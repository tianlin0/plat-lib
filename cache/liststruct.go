package cache

import (
	"container/list"
	"sync"
	"time"
)

type listStruct struct {
	l *list.List

	count int //总数量
	lock  sync.Mutex
}

func initListCache() *listStruct {
	cacheList := new(listStruct)
	cacheList.init()
	return cacheList
}

func (c *listStruct) init() {
	if c.l == nil {
		c.l = list.New()
	}
}

func (c *listStruct) push(cacheType string, key string, value interface{}, seconds time.Duration) {
	c.init()
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

	c.l.PushBack(entry)
	c.count = c.count + 1
}

func (c *listStruct) pop() (retData interface{}, canUse bool) {
	if c.l == nil {
		return nil, false
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	iter := c.l.Back()
	if iter == nil {
		return nil, false //表示删除到最后一个了
	}
	v := iter.Value
	c.l.Remove(iter)
	c.count = c.count - 1
	if val, ok := v.(*DataEntry); ok {
		return val.Value, dataValidate(val)
	}
	return nil, false
}

func (c *listStruct) get(getKey string) (interface{}, bool) {
	c.init()

	var e *DataEntry

	for iter := c.l.Front(); iter != nil; iter = iter.Next() {
		v := iter.Value
		if vTemp, ok := v.(*DataEntry); ok {
			if vTemp.Key == getKey {
				e = vTemp
				break
			}
		}
	}
	if e == nil {
		return nil, false
	}

	return e.Value, dataValidate(e)
}

func (c *listStruct) getEntry(key string) *DataEntry {
	c.init()

	var e *DataEntry

	for iter := c.l.Front(); iter != nil; iter = iter.Next() {
		v := iter.Value
		if vTemp, ok := v.(*DataEntry); ok {
			if vTemp.Key == key {
				e = vTemp
				break
			}
		}
	}
	if e == nil {
		return nil
	}

	return e
}

// 数据有效
func dataValidate(e *DataEntry) bool {
	//没到过期时间
	expired := time.Now().Before(e.Expire)
	if expired {
		return true
	}
	//如果永久在线
	isOnline := e.Expire.Before(e.CTime)
	if isOnline {
		return true
	}
	return false
}

func (c *listStruct) del(delKey string) bool {
	c.init()

	c.lock.Lock()
	defer c.lock.Unlock()

	delOk := false
	for iter := c.l.Front(); iter != nil; iter = iter.Next() {
		v := iter.Value
		if vTemp, ok := v.(*DataEntry); ok {
			if vTemp.Key == delKey {
				c.l.Remove(iter)
				delOk = true
				break
			}
		}
	}

	if delOk {
		c.count = c.count - 1
	}
	return delOk
}

// rangeBack 返回 false 则跳出循环
func (c *listStruct) ranges(rangeBack func(dataEntry *DataEntry) bool) bool {
	c.init()
	retBool := true
	for iter := c.l.Front(); iter != nil; iter = iter.Next() {
		v := iter.Value
		if vTemp, ok := v.(*DataEntry); ok {
			retBool = rangeBack(vTemp)
			if !retBool {
				break
			}
		}
	}
	return retBool
}

// 清除过期的，最大的数量maxLen，多于这个数的，则会清空时间最久的，返回目前执行后的数据数量
// deleteCallback 返回 true 表示删除了
func (c *listStruct) flush(deleteCallback func(dataEntry *DataEntry, isEmpty bool, isDelete bool) bool) int {
	var runListCount = 0

	for iter := c.l.Front(); iter != nil; iter = iter.Next() {
		runListCount++
		value := iter.Value
		e, ok := value.(*DataEntry)
		isWillDelete := false //因为过期将要删除这个选项
		if !ok {
			//格式不对
			isWillDelete = true
			e = nil
		} else {
			okData := dataValidate(e)
			if !okData {
				isWillDelete = true
			}
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
		if e == nil {
			deleteKey = true
		}

		//过期将要删除，然后手动函数也返回能删除，则该项目会被删除掉
		if deleteKey {
			c.lock.Lock()
			c.l.Remove(iter)
			c.count = c.count - 1
			runListCount--
			c.lock.Unlock()
		}
	}
	return runListCount
}
