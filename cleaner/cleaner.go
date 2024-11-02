// Package cleaner s.go
package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Cleanable 清理器
type Cleanable interface {
	Stop()
	Name() string
}

var (
	resourcesMu sync.RWMutex
	resources   = make([]Cleanable, 1)
)

// Register 注册清理器
func Register(r ...Cleanable) {
	resourcesMu.Lock()
	defer resourcesMu.Unlock()
	resources = append(resources, r...)
}

// Run 运行清理器
func Run(ctx context.Context) {
	unRegisterAll := func() {
		resourcesMu.Lock()
		defer resourcesMu.Unlock()
		resources = make([]Cleanable, 1)
	}

	var wg sync.WaitGroup
	wg.Add(len(resources))
	cleanup := func(reason string) {
		last := len(resources) - 1
		for i := range resources {
			r := resources[last-i]
			if r != nil {
				fmt.Printf("( %s ) terminated, %s", r.Name(), reason)
				r.Stop()
			}
			wg.Done()
		}
		unRegisterAll()
	}

	terminateIf(ctx,
		func() {
			cleanup("cancel")
		},
		func(s os.Signal) {
			cleanup(fmt.Sprintf("signal %+v", s))
		})
	wg.Wait()
}

type onCancel func()
type onSignal func(os.Signal)

func terminateIf(ctx context.Context, onCancel onCancel, onSignal onSignal) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGILL, syscall.SIGTERM,
		syscall.SIGTRAP, syscall.SIGQUIT, syscall.SIGABRT)
	go func() {
		for {
			select {
			case <-ctx.Done():
				onCancel()
				return
			case s := <-sig:
				onSignal(s)
				return
			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
}
