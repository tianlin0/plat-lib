package dtms

import (
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

//跨语言分布式事务管理器

func Start() {
	const qsBusi = "http://localhost:8081/api/busi_saga"
	req := &gin.H{"amount": 30} // 微服务的载荷
	// DtmServer为DTM服务的地址，是一个url
	DtmServer := "http://localhost:36789/api/dtmsvr"
	saga := dtmcli.NewSaga(DtmServer, shortuuid.New()).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 补偿操作为url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransIn"， 补偿操作为url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()

	fmt.Println(err)
}

//import (
//	"os"
//	"time"
//
//	"github.com/dtm-labs/dtm-examples/examples"
//	"github.com/dtm-labs/dtm/dtmsvr"
//	"github.com/sirupsen/logrus"
//)
//
//// M alias
//type M = map[string]interface{}
//
//func wait() {
//	for {
//		time.Sleep(10000 * time.Second)
//	}
//}
//
//// StartServer
//func StartServer() {
//	if len(os.Args) > 1 && os.Args[1] == "dtmsvr" { // 实际运行，只启动dtmsvr，不重新load数据
//		dtmsvr.StartSvr()
//		//dtmsvr.MainStart()
//		wait()
//	}
//	// 下面都是运行示例，因此首先把服务器的数据重新准备好
//	dtmsvr.PopulateMysql()
//	dtmsvr.MainStart()
//	if len(os.Args) == 1 { // 默认没有参数的情况下，准备好数据并启动dtmsvr即可
//		wait()
//	}
//	// quick_start 比较独立，单独作为一个例子运行，方便新人上手
//	if len(os.Args) > 1 && (os.Args[1] == "quick_start" || os.Args[1] == "qs") {
//		examples.QsStartSvr()
//		examples.QsFireRequest()
//		wait()
//	}
//
//	// 下面是各类的例子
//	examples.PopulateMysql()
//	app := examples.BaseAppStartup()
//	if os.Args[1] == "xa" { // 启动xa示例
//		examples.XaSetup(app)
//		examples.XaFireRequest()
//	} else if os.Args[1] == "saga" { // 启动saga示例
//		examples.SagaSetup(app)
//		examples.SagaFireRequest()
//	} else if os.Args[1] == "tcc" { // 启动tcc示例
//		examples.TccSetup(app)
//		examples.TccFireRequest()
//	} else if os.Args[1] == "msg" { // 启动msg示例
//		examples.MsgSetup(app)
//		examples.MsgFireRequest()
//	} else if os.Args[1] == "all" { // 运行所有示例
//		examples.SagaSetup(app)
//		examples.TccSetup(app)
//		examples.XaSetup(app)
//		examples.MsgSetup(app)
//		examples.SagaFireRequest()
//		examples.TccFireRequest()
//		examples.XaFireRequest()
//		examples.MsgFireRequest()
//	} else if os.Args[1] == "saga_barrier" {
//		examples.SagaBarrierAddRoute(app)
//		examples.SagaBarrierFireRequest()
//	} else if os.Args[1] == "tcc_barrier" {
//		examples.TccBarrierAddRoute(app)
//		examples.TccBarrierFireRequest()
//	} else {
//		logrus.Fatalf("unknown arg: %s", os.Args[1])
//	}
//	wait()
//}
