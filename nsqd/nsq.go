package nsqd

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/tianlin0/plat-lib/curl"
	"net/http"
)

// StartConsumer 启动一个nsq监听消费
func StartConsumer(lookUpHttpAddr string, topicName string, channelName string, handle nsq.HandlerFunc) {
	config := nsq.NewConfig()
	consumer, _ := nsq.NewConsumer(topicName, channelName, config)
	consumer.AddHandler(handle)
	err := consumer.ConnectToNSQLookupd(lookUpHttpAddr)
	if err == nil {
		select {}
	}
}

// PushMessage 发送消息，默认为：4151
func PushMessage(nsqPubOrigin string, headers http.Header, topicName string, message string) (bool, error) {
	pubUrl := fmt.Sprintf("%s/pub?topic=%s", nsqPubOrigin, topicName)
	resp := curl.NewRequest(&curl.Request{
		Url:    pubUrl,
		Data:   message,
		Method: http.MethodPost,
		Header: headers,
	}).SetRetryPolicy(&curl.RetryPolicy{
		RespDateType: "string",
	}).Submit(nil)
	if resp.Error != nil {
		return false, resp.Error
	}
	if resp.Response == "OK" {
		return true, nil
	}
	return false, nil
}
