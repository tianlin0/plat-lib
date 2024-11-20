package templates

import (
	"fmt"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
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

func TestRuleGo(t *testing.T) {
	ruleFile := []byte(`{
  "ruleChain": {
    "id":"chain_call_rest_api",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
      {
        "id": "s1",
        "type": "jsFilter",
        "name": "过滤",
        "debugMode": true,
        "configuration": {
          "jsScript": "return msg!='bb';"
        }
      },
      {
        "id": "s2",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['test']='test02';\n metadata['index']=52;\n msgType='TEST_MSG_TYPE2';\n  msg['aa']=66; return {'msg':msg,'metadata':metadata,'msgType':msgType};"
        }
      },
      {
        "id": "s3",
        "type": "restApiCall",
        "name": "推送数据",
        "debugMode": true,
        "configuration": {
          "restEndpointUrlPattern": "http://192.168.136.26:9099/api/msg",
          "requestMethod": "POST",
          "maxParallelRequestsCount": 200
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "s2",
        "type": "True"
      },
      {
        "fromId": "s2",
        "toId": "s3",
        "type": "Success"
      }
    ]
  }
}`)
	ruleEngine, _ := rulego.New("rule01", ruleFile)

	metaData := types.NewMetadata()
	metaData.PutValue("productType", "test01")

	msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, `{"temperature":35}`)

	ruleEngine.OnMsg(msg)

}
