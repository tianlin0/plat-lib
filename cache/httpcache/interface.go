package httpcache

import (
	"context"
	"github.com/tianlin0/plat-lib/internal/gocache/lib/store"
	"time"
)

type EvictionPolicy int

const (
	LRUPolicy EvictionPolicy = iota
	FIFOPolicy
	LFUPolicy
	RandomPolicy
)

// HttpCache 获取某一个数据的接口
type HttpCache[P any, V any] interface {
	Get(ctx context.Context, cacheKey string, dataParam P) (V, error)
	Set(ctx context.Context, cacheKey string, dataValue V) bool
	Del(ctx context.Context, cacheKey string) bool
}

// Config 配置
type Config[P any, V any] struct {
	Namespace              string                                                             //全局唯一，保证存储的一类数据，数据分类使用
	StoreList              []store.StoreInterface                                             //存储的类型，可以有多个，这样可以比如有内存和redis共同存储
	MaxSize                int                                                                //存储的最大数量，控制存储数量，避免内存过大
	EvictionType           EvictionPolicy                                                     //未过有效期，超过MaxSize后主动淘汰的策略
	Timeout                time.Duration                                                      //获取超时时间，有可能硬盘出现问题，存在缓存慢的情况，如果超时，则执行ExecuteGetDataHandle，默认不设置
	Expiration             time.Duration                                                      //数据多长时间过期，过期以后被动淘汰
	CleanupInterval        time.Duration                                                      //间隔过久执行清理，主动清理
	AsyncExecuteDuration   time.Duration                                                      //在这段时间里不执行异步更新，避免瞬时压力
	NeedAsyncExecuteHandle func(ctx context.Context, dataValue V) bool                        //这个数据是否需要自动异步更新
	GetDataExecuteHandle   func(ctx context.Context, cacheKey string, dataParam P) (V, error) //动态获取数据
}

type cacheIns[P any, V any] struct {
	cfg *Config[P, V]
}

type cacheData[V any] struct {
	data           V         //存储数据
	createTime     time.Time //创建时间
	expirationTime time.Time //过期时间
}
