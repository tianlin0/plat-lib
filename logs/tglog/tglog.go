// Package tglog tgLog日志
package tglog

import (
	"github.com/tianlin0/plat-lib/cache"
)

const (
	cacheTypeName = "tgLogger"
	cacheTimeout  = 24 * 3600 * 2
)

var tgLogInstance = cache.NewMapCache(&cache.MapCache{
	CacheType: cacheTypeName,
	MaxLen:    2,
	FlushCallback: func(cacheData *cache.DataEntry, isAllEmpty bool, isDelete bool) bool {
		if isAllEmpty || cacheData == nil {
			return true
		}
		if isDelete {
			if oneLogTemp, ok := cacheData.Value.(*TGLoggerV1); ok {
				err := oneLogTemp.UdpConn.Close()
				if err == nil {
					return true
				}
			}
		}
		return false
	},
})
