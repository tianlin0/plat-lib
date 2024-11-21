package lock

import (
	"fmt"
	"github.com/tianlin0/plat-lib/cache"
	"github.com/tianlin0/plat-lib/conn"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	cache.SetDefaultRedis(&conn.Connect{
		Driver:   "redis",
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: "0",
	})
	go func() {
		mm := Lock("abcde", func() {
			fmt.Println("内 11111111")
			time.Sleep(3 * time.Second)
			fmt.Println("内 222222222")
		})
		fmt.Println("内", mm)
	}()
	mm := Lock("abcde", func() {
		fmt.Println("11111111")
		time.Sleep(5 * time.Second)
		fmt.Println("222222222")
	})
	fmt.Println("外", mm)
}
