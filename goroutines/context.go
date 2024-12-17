package goroutines

import (
	"context"
	"github.com/go-eden/routine"
	"sync"
)

var ctxMap sync.Map

// SetContext 设置上下文
func SetContext(ctx *context.Context) {
	ctxMap.Store(routine.Goid(), ctx)
}

// GetContext 获取上下文
func GetContext() (ctx *context.Context) {
	val, ok := ctxMap.Load(routine.Goid())
	if !ok {
		return nil
	}
	ctx, _ = val.(*context.Context)
	return
}

// DelContext 删除上下文
func DelContext() {
	ctxMap.Delete(routine.Goid())
}
