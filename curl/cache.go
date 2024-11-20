package curl

import (
	"github.com/json-iterator/go"
	"github.com/tianlin0/plat-lib/cache"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/goroutines"
	"time"
)

// responseCacheStruct 返回的缓存结构
type responseCacheStruct struct {
	Time     time.Time `json:"time"`
	Response string    `json:"response"`
}

func getDataFromCache(p *genRequest) string {
	if p.Cache == 0 {
		return ""
	}

	cacheId := getCacheKey(p.getRequest())

	cacheNew := cache.New()
	retData, err := cacheNew.Get(cacheId)
	if err != nil || retData == "" {
		return ""
	}

	cacheData := new(responseCacheStruct)
	err = jsoniter.Unmarshal([]byte(retData), cacheData)
	if err != nil {
		cacheNew.Del(cacheId)
		return ""
	}
	//超时
	if time.Now().Sub(cacheData.Time) > p.Cache {
		cacheNew.Del(cacheId)
		return ""
	}

	findData := false
	if p.retryPolicy != nil {
		retOk, err := p.retryPolicy.onlyCheckCondition(cacheData.Response)
		if err == nil && !retOk {
			findData = true
		}
	}

	if findData {
		return cacheData.Response
	}

	return ""
}

func setDataToCache(p *Response, cacheTime time.Duration) {
	goroutines.GoAsyncHandler(func(params ...interface{}) {
		cacheData := responseCacheStruct{
			Time:     time.Now(),
			Response: p.Response,
		}
		cacheStr := conv.String(cacheData)
		if cacheStr == "" {
			return
		}

		cacheId := getCacheKey(p.Request)

		cacheNew := cache.New()
		_, _ = cacheNew.Set(cacheId, string(cacheStr), cacheTime)
	}, nil)
}
