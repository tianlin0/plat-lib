package session

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/encode"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var cookieSessionMap = new(sync.Map)

// NewCookieStore cookie初始化
func NewCookieStore(keyPairs ...string) *sessions.CookieStore {
	keyList := make([]string, 0)
	for _, one := range keyPairs {
		keyList = append(keyList, one)
	}
	//表示session最长一月生成一个
	if len(keyList) == 0 {
		keyList = append(keyList, conv.FormatFromUnixTime(conv.ShortMonthForm06))
	}

	sort.Strings(keyList)
	cacheKey := encode.Md5(strings.Join(keyList, "|"))

	pairsList := make([][]byte, 0)
	for _, one := range keyList {
		pairsList = append(pairsList, []byte(one))
	}

	storeCache, ok := cookieSessionMap.Load(cacheKey)
	if ok {
		storeTemp, ok := storeCache.(*sessions.CookieStore)
		if ok {
			return storeTemp
		}
	}

	store := sessions.NewCookieStore(pairsList...)
	cookieSessionMap.Store(cacheKey, store)
	return store
}

// Get 获取session信息
func Get(store *sessions.CookieStore, r *http.Request, sessionName string,
	keyName string) (interface{}, error) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}
	if v, ok := session.Values[keyName]; ok {
		return v, nil
	}

	return nil, nil
}

// Save 保存session信息
func Save(store *sessions.CookieStore, r *http.Request, w http.ResponseWriter,
	sessionName string, dataMap map[string]interface{}, age int) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	for key, val := range dataMap {
		session.Values[key] = val
	}

	session.Options.MaxAge = age

	err = session.Save(r, w)
	if err != nil {
		for _, val := range dataMap {
			gob.Register(val)
		}
	}
	return err
}

// DeleteAll 删除所有
func DeleteAll(store *sessions.CookieStore, r *http.Request, w http.ResponseWriter, sessionName string) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	return err
}
