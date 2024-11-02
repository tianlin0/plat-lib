package templates

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"regexp"
	"strings"
	tmpl "text/template"
)

const (
	prefixDefault = "{{"
	suffixDefault = "}}"
)

type template interface {
	Replace(replaceMap ...interface{}) (string, error)
	GetKeyList(replaceMap ...string) []string
}

type impl struct {
	prefix string
	suffix string
	s      string
}

// NewTemplate 新建一个模板
func NewTemplate(s string, fix ...string) template {
	var prefix, suffix string
	if len(fix) == 1 {
		prefix = fix[0]
	} else if len(fix) >= 2 {
		prefix = fix[0]
		suffix = fix[1]
	}
	if prefix == "" {
		prefix = prefixDefault
	}
	if suffix == "" {
		suffix = suffixDefault
	}

	return &impl{prefix: prefix, suffix: suffix, s: s}
}

// Replace 替换对象的值
func (t *impl) Replace(replaceMap ...interface{}) (string, error) {
	tempStr := t.s
	for _, one := range replaceMap {
		newCurrParam := conv.ToKeyListFromMap(one)
		for key, val := range newCurrParam {
			valStr := conv.String(val)
			if key == "." { //表示填充所有字符串的情况，用来做特殊处理的
				tempStr = strings.ReplaceAll(tempStr, t.prefix+t.suffix, valStr)
				continue
			}
			tempStr = strings.ReplaceAll(tempStr, t.prefix+key+t.suffix, valStr)
		}
	}
	return tempStr, nil
}

// GetKeyList GetTemplateKeys 检查
func (t *impl) GetKeyList(replaceMap ...string) []string {
	allList := make([]string, 0)
	for _, one := range replaceMap {
		reg := regexp.MustCompile(t.prefix + "(.*?)" + t.suffix)
		//返回str中第一个匹配reg的字符串
		data := reg.FindAllString(one, -1)
		if data != nil && len(data) > 0 {
			for _, temp := range data {
				temp = strings.ReplaceAll(temp, t.prefix, "")
				temp = strings.ReplaceAll(temp, t.suffix, "")
				allList = append(allList, temp)
			}
		}
	}
	return allList
}

// Template 模版填充
/*
定义变量：{{$article := "hello"}}  {{$article := .ArticleContent}}
调用方法：{{functionName .arg1 .arg2}}
条件判断：
{{if .condition1}}
{{else if .condition2}}
{{end}}

逻辑关系：
or and not : 或, 与, 非
eq ne lt le gt ge: 等于, 不等于, 小于, 小于等于, 大于, 大于等于
示例：
{{if ge .var1 .var2}}
{{end}}

循环：
{{range $i, $v := .slice}}
{{end}}

{{range .slice}}
{{.field}} //获取对象里的变量
{{$.ArticleContent}}  //访问循环外的全局变量的方式
{{end}}
*/
func Template(format string, data interface{}, delis ...string) (string, error) {
	startDelis := prefixDefault
	endDelis := suffixDefault

	if len(delis) >= 2 {
		if delis[0] != "" {
			startDelis = delis[0]
		}
		if delis[1] != "" {
			endDelis = delis[1]
		}
	} else if len(delis) == 1 {
		if delis[0] != "" {
			startDelis = delis[0]
		}
	}

	{ //key里面有/将执行报错，需要处理一下/的问题

	}

	tempName := toMd5(format)
	var err error
	templateIns := tmpl.New("utils-template-"+tempName).Delims(startDelis, endDelis)
	templateIns = templateIns.Option("missingkey=zero")
	templateIns, err = templateIns.Parse(format)
	if err != nil {
		return "", err
	}
	var doc bytes.Buffer
	err = templateIns.Execute(&doc, data)
	if err == nil {
		// 如果没有值也不会报错，所以这里需要处理一下
		docStr := doc.String()
		HasNoValueIndex := strings.Index(docStr, "<no value>")
		if HasNoValueIndex < 0 {
			return docStr, nil
		}
		return docStr, fmt.Errorf("模板有未覆盖的变量")
	}

	return doc.String(), err
}

func toMd5(s string) string {
	d := md5.Sum([]byte(s))
	return hex.EncodeToString(d[:])
}
