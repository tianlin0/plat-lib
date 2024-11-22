package database

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/tianlin0/plat-lib/logs"
	"reflect"
	"strings"
	"testing"
)

func TestBachGetPageList(t *testing.T) {
	//aaa, bbb := GetColumnLikeSql("ddddfdfdsfsd2323_%%")
	//fmt.Println(aaa, bbb)

	oneDataCacheMap := cmap.New()
	CacheType := "aaaa"
	dataOld := "aa"

	oneDataCacheMap.Set(CacheType, reflect.TypeOf(dataOld))

	dataTemp := "bbb"
	if value, ok := oneDataCacheMap.Get(CacheType); ok {
		if reflect.TypeOf(dataTemp) != value {
			logs.DefaultLogger().Error("executeHandle data type error:", CacheType, dataTemp)
		}
	}
}
func TestBachGetPageList1(t *testing.T) {
	file := ""
	currentLib := ""

	aa := strings.Index(file, currentLib)

	fmt.Println(aa)
}
