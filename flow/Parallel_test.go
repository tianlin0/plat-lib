package flow_test

import (
	"fmt"
	"github.com/tianlin0/plat-lib/flow"
	"sync/atomic"
	"testing"
)

func TestParallel(t *testing.T) {
	a := int64(0)
	b := int64(0)

	flow.Parallel(func() {
		atomic.AddInt64(&a, 1)
	}, func() {
		atomic.AddInt64(&b, 1)
	})

	fmt.Println(a, b)
}

func TestParallelRepeat(t *testing.T) {
	a := int64(0)
	b := int64(0)
	n := 10

	flow.ParallelRepeat(n, func() {
		atomic.AddInt64(&a, 1)
	}, func() {
		atomic.AddInt64(&b, 1)
	})

	fmt.Println(a, b)
}
