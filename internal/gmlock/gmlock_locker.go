// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock

import (
	gocache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// Locker is a memory based locker.
// Note that there's no cache expire mechanism for mutex in locker.
// You need remove certain mutex manually when you do not want use it anymore.
type Locker struct {
	m *gocache.Cache
	s *gocache.Cache
}

var (
	expiration      = 10 * time.Minute
	cleanupInterval = 30 * time.Minute
)

// New creates and returns a new memory locker.
// A memory locker can lock/unlock with dynamic string key.
func New() *Locker {
	return &Locker{
		m: gocache.New(expiration, cleanupInterval),
		s: gocache.New(expiration, cleanupInterval),
	}
}

// Lock locks the `key` with writing lock.
// If there's a write/reading lock the `key`,
// it will block until the lock is released.
func (l *Locker) Lock(key string) {
	l.s.Set(key, true, 0)
	l.getOrNewMutex(key).Lock()
	l.s.Set(key, true, 0)
}

// TryLock tries locking the `key` with writing lock,
// it returns true if success, or it returns false if there's a writing/reading lock the `key`.
func (l *Locker) TryLock(key string) bool {
	l.s.Set(key, true, 0)
	return l.getOrNewMutex(key).TryLock()
}

// Unlock unlocks the writing lock of the `key`.
func (l *Locker) Unlock(key string) {
	if v, ok := l.m.Get(key); v != nil && ok {
		l.s.Set(key, false, 0)
		v.(*sync.RWMutex).Unlock()
	}
}

// RLock locks the `key` with reading lock.
// If there's a writing lock on `key`,
// it will blocks until the writing lock is released.
func (l *Locker) RLock(key string) {
	l.s.Set(key, true, 0)
	l.getOrNewMutex(key).RLock()
}

// TryRLock tries locking the `key` with reading lock.
// It returns true if success, or if there's a writing lock on `key`, it returns false.
func (l *Locker) TryRLock(key string) bool {
	l.s.Set(key, true, 0)
	return l.getOrNewMutex(key).TryRLock()
}

// RUnlock unlocks the reading lock of the `key`.
func (l *Locker) RUnlock(key string) {
	if v, ok := l.m.Get(key); v != nil && ok {
		l.s.Set(key, false, 0)
		v.(*sync.RWMutex).RUnlock()
	}
}

// LockFunc locks the `key` with writing lock and callback function `f`.
// If there's a write/reading lock the `key`,
// it will block until the lock is released.
//
// It releases the lock after `f` is executed.
func (l *Locker) LockFunc(key string, f func()) {
	l.Lock(key)
	defer l.Unlock(key)
	f()
}

// RLockFunc locks the `key` with reading lock and callback function `f`.
// If there's a writing lock the `key`,
// it will block until the lock is released.
//
// It releases the lock after `f` is executed.
func (l *Locker) RLockFunc(key string, f func()) {
	l.RLock(key)
	defer l.RUnlock(key)
	f()
}

// TryLockFunc locks the `key` with writing lock and callback function `f`.
// It returns true if success, or else if there's a write/reading lock the `key`, it return false.
//
// It releases the lock after `f` is executed.
func (l *Locker) TryLockFunc(key string, f func()) bool {
	if l.TryLock(key) {
		defer l.Unlock(key)
		f()
		return true
	}
	return false
}

// TryRLockFunc locks the `key` with reading lock and callback function `f`.
// It returns true if success, or else if there's a writing lock the `key`, it returns false.
//
// It releases the lock after `f` is executed.
func (l *Locker) TryRLockFunc(key string, f func()) bool {
	if l.TryRLock(key) {
		defer l.RUnlock(key)
		f()
		return true
	}
	return false
}

// Remove removes mutex with given `key` from locker.
func (l *Locker) Remove(key string) bool {
	// 解锁状态下才可以删除
	locked := l.Locked(key)
	if locked {
		return false
	}
	l.m.Delete(key)
	l.s.Delete(key)
	return true
}

// Locked removes mutex with given `key` from locker.
func (l *Locker) Locked(key string) bool {
	b, ok := l.s.Get(key)
	if ok {
		bT, ok := b.(bool)
		if ok {
			return bT
		}
	}
	return false
}

// Clear removes all mutexes from locker.
func (l *Locker) Clear() {
	l.m.Flush()
	l.s.Flush()
}

// getOrNewMutex returns the mutex of given `key` if it exists,
// or else creates and returns a new one.
func (l *Locker) getOrNewMutex(key string) *sync.RWMutex {
	if v, ok := l.m.Get(key); v != nil && ok {
		if k, ok := v.(*sync.RWMutex); ok {
			l.m.Set(key, k, 0)
			return k
		}
	}
	newMutex := new(sync.RWMutex)
	l.m.Set(key, newMutex, 0)

	return newMutex
}
