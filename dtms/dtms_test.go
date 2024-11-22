package dtms

import (
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"testing"
)

//跨语言分布式事务管理器

func TestAesCbC(t *testing.T) {

	const qsBusi = "http://localhost:8081/api/busi_saga"
	req := &gin.H{"amount": 30} // 微服务的载荷
	// DtmServer为DTM服务的地址，是一个url
	DtmServer := "http://localhost:36789/api/dtmsvr"
	saga := dtmcli.NewSaga(DtmServer, uuid.New().String()).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 补偿操作为url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransIn"， 补偿操作为url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()

	fmt.Println(err)
}
