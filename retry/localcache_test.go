package retry

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type AA struct {
	Name string
}

func TestCacheMap(t *testing.T) {

	str := "Hell!"

	// 获取前3个字符
	firstThree := str[:6]
	fmt.Println(firstThree)

	var a AA
	err := New().WithInterval(1*time.Second).WithAttemptCount(7).Do(nil, func(ctx context.Context) (interface{}, error) {
		return a, nil
	}, &a)

	fmt.Println(err, a)
}
