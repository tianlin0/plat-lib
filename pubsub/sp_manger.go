package pubsub

import (
	"sync"
	"time"
)

//DELSP_TICKER 定时器删除空订阅的间隔 秒
const DELSP_TICKER = 3

type spManger struct {
	spDic map[interface{}]*pubSub
	sync.RWMutex
}

func GetEvent(i interface{}) *pubSub {
	sp, ok := manger.getSp(i)
	if !ok {
		sp = newPubSub()
		manger.addSp(i, sp)
	}
	return sp
}

var manger *spManger

func init() {
	manger = &spManger{
		spDic: make(map[interface{}]*pubSub, 0),
	}
	go manger.startTicker()
}

func (m *spManger) getSp(i interface{}) (*pubSub, bool) {
	manger.Lock()
	defer manger.Unlock()
	sp, ok := manger.spDic[i]
	return sp, ok

}

func (m *spManger) addSp(i interface{}, sp *pubSub) {
	m.Lock()
	defer m.Unlock()
	m.spDic[i] = sp
}

func (m *spManger) delSp(sp *pubSub) {
	m.Lock()
	defer m.Unlock()
	delete(m.spDic, sp)
}

//开启一个定时器检查是否有数据
func (m *spManger) startTicker() {
	ticker := time.NewTicker(time.Second * DELSP_TICKER)
	for range ticker.C {
		for _, sp := range m.spDic {
			if sp.getSubscribesCount() == 0 {
				sp.cancel()
				m.delSp(sp)
			}
		}
	}
}
