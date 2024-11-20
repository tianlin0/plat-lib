package goroutines

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// GoSyncHandler 同步方法
func GoSyncHandler(callFunction func(params ...interface{}), panicHandle func(err error), params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			//打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			stackInfo = strings.ReplaceAll(stackInfo, "\n", "|")
			errStr := fmt.Sprintf("panic_stack_info: %s ### %s", err, stackInfo)
			log.Println(errStr)
			if panicHandle != nil {
				panicHandle(fmt.Errorf(errStr))
			}
			return
		}
	}()
	callFunction(params...)
}

// GoAsyncHandler 异步方法
// 异步有一个总量，不然会创建太多，造成系统阻塞了
// var ch = make(chan bool, 50)
func GoAsyncHandler(callFunction func(params ...interface{}), panicHandle func(err error), params ...interface{}) {
	//ch <- true
	go func(tempParams ...interface{}) {
		//defer func() {
		//	<-ch
		//}()
		GoSyncHandler(callFunction, panicHandle, tempParams...)
	}(params...)
}

// 并行执行最大效率
func setRunTimeProcess() bool {
	pocsSet := os.Getenv("GOMAXPROCS")
	// 环境变量已设置，就不用设置了。
	if pocsSet == "" {
		cpuNumStr := os.Getenv("CPU_LIMIT")
		cpuNum, _ := strconv.ParseFloat(cpuNumStr, 64)
		cpuNumInt := int(math.Ceil(cpuNum))
		if cpuNumInt <= 0 || cpuNumInt >= 20 {
			cpuNumInt = runtime.NumCPU()
		}

		// 如果是宿主机的数量,太大的话，则默认改为10
		if cpuNumInt <= 0 || cpuNumInt >= 20 {
			cpuNumInt = 10
		}
		runtime.GOMAXPROCS(cpuNumInt)
		return true
	}
	return false
}

func init() {
	setRunTimeProcess()
}
