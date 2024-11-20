package utils

import (
	"sort"
	"sync"
)

// NewSortSet 创建协程安全的HashSet
func NewSortSet() *syncHashSet {
	return &syncHashSet{
		values: hashSortSet{},
	}
}

// 协程安全的HashSet
type syncHashSet struct {
	values  hashSortSet
	sortNum int // 按顺序输出
	mutex   sync.Mutex
}

// Add 往set添加元素
func (s *syncHashSet) Add(value interface{}) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sortNum = s.sortNum + 1
	return s.values.add(value, s.sortNum)
}

// Delete 删除元素
func (s *syncHashSet) Delete(value interface{}) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.values.delete(value)
}

// Contains 检查元素存在性
func (s *syncHashSet) Contains(value interface{}) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.values.contains(value)
}

// List 返回列表的元素
func (s *syncHashSet) List() []interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	list := s.values.copy()

	listTemp := make([]interface{}, 0)

	listKey := make([]int, 0)
	listOld := make(map[int]interface{})
	for one, sortTemp := range list {
		temp1 := sortTemp
		temp2 := one
		listKey = append(listKey, temp1)
		listOld[temp1] = temp2
	}
	sort.Ints(listKey)
	for _, one := range listKey {
		if temp, ok := listOld[one]; ok {
			listTemp = append(listTemp, temp)
		}
	}
	return listTemp
}

// Hash集合数据结构
type hashSortSet map[interface{}]int

// add 往集合添加值
func (h hashSortSet) add(value interface{}, sortNum int) bool {
	if _, ok := h[value]; !ok {
		h[value] = sortNum
		return true
	}
	return false
}

// delete 往集合删除值
func (h hashSortSet) delete(value interface{}) bool {
	if _, ok := h[value]; ok {
		delete(h, value)
		return true
	}
	return false
}

// contains 值是否存在集合中
func (h hashSortSet) contains(value interface{}) bool {
	_, ok := h[value]
	return ok
}

// copy 复制hashSet
func (h hashSortSet) copy() hashSortSet {
	newSet := make(map[interface{}]int, len(h))
	for k, v := range h {
		newSet[k] = v
	}
	return newSet
}
