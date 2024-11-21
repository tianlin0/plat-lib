package pubsub

import (
	"context"
	"sync"
)

type pubSub struct {
	subscribes []func(interface{})
	channel    chan interface{}
	ctx        context.Context
	cancel     context.CancelFunc
	sync.RWMutex
}

//开启协程来轮巡该订阅是否有数据传入
func (sp *pubSub) startSubscribe() {
	for {
		select {
		case <-sp.ctx.Done():
			break
		case obj := <-sp.channel:
			for _, fun := range sp.subscribes {
				fun(obj)
			}
		}
	}
}

func newPubSub() *pubSub {
	ctx, cancel := context.WithCancel(context.Background())
	sp := &pubSub{
		ctx:     ctx,
		cancel:  cancel,
		channel: make(chan interface{}),
	}
	go sp.startSubscribe()
	return sp
}

func (sp *pubSub) Publish(i interface{}) {
	go func() {
		sp.channel <- i
	}()
}

func (sp *pubSub) Subscribe(fun func(interface{})) {
	sp.Lock()
	defer sp.Unlock()
	sp.subscribes = append(sp.subscribes, fun)
}

func (sp *pubSub) getSubscribesCount() int {
	sp.Lock()
	defer sp.Unlock()
	return len(sp.subscribes)
}
