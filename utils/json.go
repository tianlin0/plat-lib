package utils

import (
	"encoding/json"
	"sort"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/tianlin0/plat-lib/conv"
)

// DiffJsonChanges 比较两json的差异
func DiffJsonChanges(old, new interface{}) (string, error) {
	a, _ := json.Marshal(old)
	b, _ := json.Marshal(new)
	patch, err := jsonpatch.CreateMergePatch(a, b)
	if err != nil {
		return "", err
	}
	return string(patch), nil
}

// GetJsonOnlyKey 传入map对象，取得唯一的返回key，用于cache中存储的时候
func GetJsonOnlyKey(data interface{}) string {
	jsonData := conv.String(data)
	jsonMap := conv.ToKeyListFromMap(jsonData)
	if len(jsonMap) > 0 {
		return conv.HttpBuildQuery(jsonMap)
	}
	strList := strings.Split(jsonData, "")
	sort.Strings(strList)
	return strings.Join(strList, "")
}
