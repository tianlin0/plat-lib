package crontab

import (
	"github.com/tianlin0/plat-lib/cache"
	"github.com/tianlin0/plat-lib/conn"
	"testing"
)
import "fmt"

func TestCrontab(t *testing.T) {

	StartCrontabJobs(0, map[string]func(){
		"*/2 * * * *": func() {
			fmt.Println("1分钟1")
		},
		"0 02 17 * *": func() {
			fmt.Println("定点1")
		},
	})

	select {}
}
func TestCrontabLockKey(t *testing.T) {
	cache.SetDefaultRedis(&conn.Connect{
		Driver: "redis",
		Host:   "127.0.0.1",
		Port:   "6379",
	})
	StartCrontabJobsByLockKey("aaa", map[string]func(){
		"* * * * *": func() {
			fmt.Println("1分钟1")
		},
	})
	StartCrontabJobsByLockKey("aaa", map[string]func(){
		"* * * * *": func() {
			fmt.Println("1分钟2")
		},
	})

	select {}
}
