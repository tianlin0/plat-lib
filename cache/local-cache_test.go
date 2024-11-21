package cache

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestCacheMap(t *testing.T) {
	var localCaches = NewMapCache(&MapCache{
		MaxLen:            5,
		FlushTimeInterval: 2 * time.Second,
	})
	localCaches.Set("aaaa", "bbbbb", 10)
	localCaches.Set("bbbb", "ccc", 15)
	localCaches.Set("cccc", "dddd", 20)
	localCaches.Set("dddd", "eeee", 25)

	localCaches.Del("aaaa")

	for {
		localCaches.Range(func(dataEntry *DataEntry) bool {
			fmt.Printf("%s", dataEntry.Key)
			return true
		})
		time.Sleep(time.Second * time.Duration(1))
	}

	time.Sleep(time.Second * time.Duration(60))

	fmt.Println("localCaches:")

}
func TestCacheMap11(t *testing.T) {

	aa := runtime.GOMAXPROCS(0)
	fmt.Print(aa)

	//printMemUsage()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
