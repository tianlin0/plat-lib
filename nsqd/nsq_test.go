// Package nsq is a generated Go source code.
package nsqd_test

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/tianlin0/plat-lib/nsqd"
	"testing"
)

func TestPushMsg(t *testing.T) {
	message, err := nsqd.PushMessage("", nil, "helloTitle", "dsfdsfsdfs")
	if err != nil {
		return
	}

	fmt.Println(message)
}

func TestCustomer(t *testing.T) {
	nsqd.StartConsumer("", "helloTitle", "chenel1", func(message *nsq.Message) error {
		fmt.Println(message)
		return nil
	})
	msgIDGood := nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'}
	msg := nsq.NewMessage(msgIDGood, nil)
	msg.DisableAutoResponse()
}
