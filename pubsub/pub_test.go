package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func TestDemo(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	GetEvent("hello").Subscribe(helloworld)

	GetEvent("hello").Publish(&user{Name: "tom", Age: 18})

	time.Sleep(2 * time.Second)

	GetEvent("hello").Publish(&user{Name: "tom", Age: 28})
}

func helloworld(i interface{}) {
	fmt.Printf("receive : %+v\n", i)
}
