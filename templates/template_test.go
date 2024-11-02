package templates

import (
	"fmt"
	"strings"
	"testing"
)

func TestBachGetPageList(t *testing.T) {
	str := "111/<no value>"
	ret := strings.Index(str, "<no value>")
	if ret < 0 {
		fmt.Println("没有")
		return
	}
	fmt.Println("有")

	postUrl, err := Template("{{.aaa/aaa}}", map[string]interface{}{
		"aaa/aaa": "55555",
	})

	fmt.Println(postUrl, err)
}
func TestRuleExp(t *testing.T) {
	str := "code.numMap.num==code.numMap.num2"
	//dataMap := map[string]interface{}{
	//	"code": map[string]interface{}{
	//		"data": true,
	//	},
	//}
	//dataMap := map[string]interface{}{
	//	"code": map[string]map[string]float32{
	//		"numMap": map[string]float32{
	//			"num":  4.5,
	//			"num2": 4.51,
	//		},
	//	},
	//}

	str = "==code"
	dataMap2 := "code"
	retOk, err := RuleExpr(str, dataMap2)
	fmt.Println(retOk, err)
}
