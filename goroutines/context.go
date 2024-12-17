package goroutines

import (
	"context"
	"github.com/go-eden/routine"
	gocache "github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

var (
	expiration      = 20 * time.Minute
	cleanupInterval = 30 * time.Minute
	baseInt         = 10
)

var ctxCache *gocache.Cache

func getCache() *gocache.Cache {
	if ctxCache == nil {
		ctxCache = gocache.New(expiration, cleanupInterval)
	}
	return ctxCache
}
func getGoId() string {
	return strconv.FormatInt(routine.Goid(), baseInt)
}

// SetContext 设置上下文
func SetContext(ctx *context.Context) {
	ctxFactory := getCache()
	ctxFactory.Set(getGoId(), ctx, gocache.DefaultExpiration)
}

// GetContext 获取上下文
func GetContext() (ctx *context.Context) {
	ctxFactory := getCache()

	val, ok := ctxFactory.Get(getGoId())
	if !ok {
		return nil
	}
	ctx, ok = val.(*context.Context)
	if ok {
		return ctx
	}
	return nil
}

// DelContext 删除上下文
func DelContext() {
	ctxFactory := getCache()
	ctxFactory.Delete(getGoId())
}
