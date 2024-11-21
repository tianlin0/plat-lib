package curl

import (
	"fmt"
	"github.com/tianlin0/plat-lib/utils"
	"net/http"
	"testing"
	"time"
)

func TestCurls(t *testing.T) {
	_ = NewRequest(&Request{
		Url: "",
		Data: map[string]interface{}{
			"aaaa": 1,
		},
		Method: "GET",
	}).SetDefaultPrintInt(-1).Submit(nil)
}

func TestSubmitDemo(t *testing.T) {
	realUrl := ``

	postData := map[string]interface{}{
		"gitProjectName": "",
	}

	filePath := ""

	fileMap, err := utils.GetAllFileContent(filePath)
	if err != nil {
		return
	}
	postData["initFiles"] = fileMap

	//TODO 文件太大，会挂掉
	curlRetStruct := NewRequest(&Request{
		Url:  realUrl,
		Data: postData,
		Header: http.Header{
			"X-Gdp-Jwt-Assertion": []string{"eyJhbGciOiJSUzI1NiIsImtpZCI6IjFiNTdjMmNmLThkZWYtNGZiZi1iZjgxLWQwYTJlMDgwMTlmYSIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJvZHAtZXh0ZXJuYWwiLCJjbGllbnRfdHlwZSI6ImFkbWluIiwiZXhwIjoxNzA0MjMyMzQ5LCJpYXQiOjE3MDQxODkxNDksInNjb3BlIjoiZ2RwYWRtaW4ifQ.W-0ZmrcjQt0yi4e2ViXciohwO5GxshfD3AA171JDtWn8bmBZy6UnA2zL1y2JipCHpOY7dmb1i0vBeXXUqBaRQPv1FOXHQpaBZEp1l9EkOaeoxvPabw7ZxryetVFqMzsIgrz9KteR5M6bPa04ZqpwkgPXRIKJGJy-NoN4Lt1yrI1HnGFBA3ZcdN3znJRxPcofgR7GzevdjzpoZqoMpHiWvtafjumq6RMuVGxwnoyeuqJhXxVQ6ci8AtkriT91xPaH9ouJmWInyDaibVXju8FyCREXDDTyOEs4muw-MDajdlZxEQOFYL8APvrxRjLlTfzLd6s-2C88rqwwrzbkMwl4Xg"},
			"X-Gdp-Userid":        []string{""},
		},
		Cache:   0,
		Timeout: 2 * time.Minute,
		Method:  http.MethodPost,
	}).Submit(nil)

	fmt.Println(curlRetStruct)

}
