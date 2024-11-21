package crontab

import (
	"github.com/robfig/cron/v3"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/lock"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"math/big"
	"sync"
	"time"
)

type cronInstance struct {
	isStart bool
	c       *cron.Cron
}

var (
	oneCrontab *cronInstance
	runningMu  sync.Mutex
	ronce      sync.Once
)

// getCron 初始化
func getCron() *cronInstance {
	ronce.Do(func() {
		oneCrontab = &cronInstance{
			isStart: false,
			c:       cron.New(),
		}
	})
	return oneCrontab
}

/*
//按分钟开始定时
crontab.StartCrontabJobs(map[string]func(){
		"* * * * *" : func(){
},
})
minute     = field(fields[1], minutes)
hour       = field(fields[2], hours)
dayOfMonth = field(fields[3], dom)
month      = field(fields[4], months)
dayOfWeek  = field(fields[5], dow)
*/

//多个服务器的问题，会很快执行，需要有个判断的标示

// StartCrontabJobsByLockKey 启动或加入定时任务，多个是避免如果时间重复的话
// lockOnlyKey 可利用redis来进行多部署环境的唯一锁定，避免多次重复执行
func StartCrontabJobsByLockKey(lockKey string, jobs ...map[string]func()) {
	if lockKey != "" {
		//表示执行时，需要进行lock的判断
		for _, job := range jobs {
			for key, oneFun := range job {
				job[key] = func() {
					lock.Lock(lockKey, oneFun)
				}
			}
		}
	}

	StartCrontabJobs(0, jobs...)
}

// StartCrontabJobs 启动定时任务，格式：分钟 小时 天 月 星期
func StartCrontabJobs(randomSleepMaxTime time.Duration, jobs ...map[string]func()) {
	runningMu.Lock()
	defer runningMu.Unlock()

	oneCron := getCron()
	if len(jobs) == 0 {
		return
	}

	loggers := logs.DefaultLogger()
	allKey := make([]string, 0)
	for _, jobMap := range jobs {
		for key, _ := range jobMap {
			allKey = append(allKey, key)
		}
	}

	result := new(big.Int)
	randomSleepBigInt := big.NewInt(int64(randomSleepMaxTime))
	secondBigInt := big.NewInt(int64(time.Second))
	result.Div(randomSleepBigInt, secondBigInt)
	randomSleepSecond := int(result.Int64())

	loggers.Info("[crontab] StartCrontabJobs start:", randomSleepSecond, allKey)

	hasSuccess := false //如果全部出错的，则不用启动
	for _, jobMap := range jobs {
		for times := range jobMap {
			var err error
			//列表里需要将所有的内容保存一份，这样就可以到时候进行删除了
			if randomSleepSecond == 0 {
				_, err = oneCron.c.AddFunc(times, jobMap[times])
			} else {
				if randomSleepSecond <= 2 {
					randomSleepSecond = 5
				}

				temp := times //必须重新设置
				_, err = oneCron.c.AddFunc(times, func() {
					// 默认设置一个最小值为2，避免生成随机为0的情况
					randomSecond := utils.Random(2, randomSleepSecond)
					time.Sleep(time.Duration(randomSecond) * time.Second)
					jobMap[temp]()
				})
			}
			if err != nil {
				loggers.Error(" [crontab] StartCrontabJobs error:", times, err.Error())
			} else {
				hasSuccess = true
			}
		}
	}

	//没有添加成功，则不用启动
	if !hasSuccess {
		return
	}

	if oneCron.isStart {
		return
	}
	oneCron.isStart = true
	//异步启动
	goroutines.GoAsyncHandler(func(params ...interface{}) {
		oneCron.c.Run()
		select {}
	}, nil)
}
