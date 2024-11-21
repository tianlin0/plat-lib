package bigcache

import (
	"context"
	"errors"
	"fmt"
	"github.com/tianlin0/plat-lib/internal/gocache/lib/store"
	"strings"
	"time"
)

// ClientInterface represents a allegro/bigcache client
type ClientInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
	Delete(key string) error
	Reset() error
}

const (
	BigcacheType       = "bigcache"
	BigcacheTagPattern = "gocache_tag_%s"
)

type storeStruct struct {
	client  ClientInterface
	options *store.Options
}

// New creates a new store to BigCache instance(s)
func New(client ClientInterface, options ...store.Option) *storeStruct {
	return &storeStruct{
		client:  client,
		options: store.ApplyOptions(options...),
	}
}

// Get returns data stored from a given key
func (s *storeStruct) Get(_ context.Context, key any) (any, error) {
	item, err := s.client.Get(key.(string))
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, store.NotFoundWithCause(errors.New("unable to retrieve data from bigcache"))
	}

	return item, err
}

// GetWithTTL Not implemented for BigcacheStore
func (s *storeStruct) GetWithTTL(_ context.Context, _ any) (any, time.Duration, error) {
	return nil, 0, errors.New("method not implemented for codec, use Get() instead")
}

// Set defines data in Bigcache for given key identifier
func (s *storeStruct) Set(ctx context.Context, key any, value any, options ...store.Option) error {
	opts := store.ApplyOptionsWithDefault(s.options, options...)

	var val []byte
	switch v := value.(type) {
	case string:
		val = []byte(v)
	case []byte:
		val = v
	default:
		return errors.New("value type not supported by Bigcache store")
	}

	err := s.client.Set(key.(string), val)
	if err != nil {
		return err
	}

	if tags := opts.Tags; len(tags) > 0 {
		s.setTags(ctx, key, tags)
	}

	return nil
}

func (s *storeStruct) setTags(ctx context.Context, key any, tags []string) {
	for _, tag := range tags {
		tagKey := fmt.Sprintf(BigcacheTagPattern, tag)
		cacheKeys := make([]string, 0)

		if result, err := s.Get(ctx, tagKey); err == nil {
			if bytes, ok := result.([]byte); ok {
				cacheKeys = strings.Split(string(bytes), ",")
			}
		}

		alreadyInserted := false
		for _, cacheKey := range cacheKeys {
			if cacheKey == key.(string) {
				alreadyInserted = true
				break
			}
		}

		if !alreadyInserted {
			cacheKeys = append(cacheKeys, key.(string))
		}

		s.Set(ctx, tagKey, []byte(strings.Join(cacheKeys, ",")), store.WithExpiration(720*time.Hour))
	}
}

// Delete removes data from Bigcache for given key identifier
func (s *storeStruct) Delete(_ context.Context, key any) error {
	return s.client.Delete(key.(string))
}

// Invalidate invalidates some cache data in Bigcache for given options
func (s *storeStruct) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	opts := store.ApplyInvalidateOptions(options...)

	if tags := opts.Tags; len(tags) > 0 {
		for _, tag := range tags {
			tagKey := fmt.Sprintf(BigcacheTagPattern, tag)
			result, err := s.Get(ctx, tagKey)
			if err != nil {
				return nil
			}

			cacheKeys := []string{}
			if bytes, ok := result.([]byte); ok {
				cacheKeys = strings.Split(string(bytes), ",")
			}

			for _, cacheKey := range cacheKeys {
				s.Delete(ctx, cacheKey)
			}
		}
	}

	return nil
}

// Clear resets all data in the store
func (s *storeStruct) Clear(_ context.Context) error {
	return s.client.Reset()
}

// GetType returns the store type
func (s *storeStruct) GetType() string {
	return BigcacheType
}
