package lrucache

import (
	"context"
	"errors"
	"fmt"
	"github.com/tianlin0/plat-lib/internal/gocache/lib/store"
	"time"
)

const (
	// LRUCacheType represents the storage type as a string value
	LRUCacheType = "lru-cache"
)

// ClientInterface https://github.com/hashicorp/golang-lru
type ClientInterface[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Add(key K, value V) (evicted bool)
	Remove(key K) bool
	Purge()
}

// LRUCacheStore is a store for GoCache (memory) library
type LRUCacheStore[K comparable, V any] struct {
	client  ClientInterface[K, V]
	options *store.Options
}

// NewLRUCache creates a new store to GoCache (memory) library instance
func NewLRUCache[K comparable, V any](client ClientInterface[K, V], options ...store.Option) *LRUCacheStore[K, V] {
	return &LRUCacheStore[K, V]{
		client:  client,
		options: store.ApplyOptions(options...),
	}
}

// Get returns data stored from a given key
func (s *LRUCacheStore[K, V]) Get(_ context.Context, key any) (any, error) {
	var err error
	if k, ok := key.(K); ok {
		value, exists := s.client.Get(k)
		if !exists {
			err = store.NotFoundWithCause(errors.New("value not found in GoCache store"))
		}
		return value, err
	}
	return nil, nil
}

// GetWithTTL returns data stored from a given key and its corresponding TTL
func (s *LRUCacheStore[K, V]) GetWithTTL(_ context.Context, key any) (any, time.Duration, error) {
	return nil, 0, nil
}

// Set defines data in GoCache memoey cache for given key identifier
func (s *LRUCacheStore[K, V]) Set(_ context.Context, key any, value any, options ...store.Option) error {
	if k, ok := key.(K); ok {
		if v, ok := value.(V); ok {
			s.client.Add(k, v)
			return nil
		}
	}
	return fmt.Errorf("key,value error")
}

// Delete removes data in GoCache memoey cache for given key identifier
func (s *LRUCacheStore[K, V]) Delete(_ context.Context, key any) error {
	if k, ok := key.(K); ok {
		s.client.Remove(k)
		return nil
	}

	return fmt.Errorf("key,value error")
}

// Invalidate invalidates some cache data in GoCache memoey cache for given options
func (s *LRUCacheStore[K, V]) Invalidate(_ context.Context, options ...store.InvalidateOption) error {
	return nil
}

// Clear resets all data in the store
func (s *LRUCacheStore[K, V]) Clear(_ context.Context) error {
	s.client.Purge()
	return nil
}

// GetType returns the store type
func (s *LRUCacheStore[K, V]) GetType() string {
	return LRUCacheType
}
