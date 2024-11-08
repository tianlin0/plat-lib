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
	resources   = make([]Cleanable, 0)
)

// Register 注册清理器
func Register(r ...Cleanable) {
	resourcesMu.Lock()
	defer resourcesMu.Unlock()
	resources = append(resources, r...)
}

// Run 运行清理器
func Run(ctx context.Context) {
	unRegisterList := func(start, end int) {
		resourcesMu.Lock()
		defer resourcesMu.Unlock()
		if start <= 0 {
			start = 0
		}
		if end >= len(resources)-1 {
			end = len(resources) - 1
		}
		if start > end {
			return
		}

		if end == len(resources)-1 {
			resources = resources[:start]
		} else {
			resources = append(resources[:start], resources[end+1:]...)
		}
	}

	var wg sync.WaitGroup
	total := len(resources)
	if total == 0 {
		return
	}
	start := 0
	end := total - 1

	wg.Add(total)
	cleanup := func(reason string) {
		for i := range resources {
			if i > end {
				break
			}
			r := resources[i]
			if r != nil {
				fmt.Printf("( %s ) terminated, %s", r.Name(), reason)
				r.Stop()
			}
			wg.Done()
		}
		unRegisterList(start, end)
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
