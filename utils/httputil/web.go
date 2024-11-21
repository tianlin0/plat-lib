package httputil

import (
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/utils"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetDomainHost 取得域名或上一级域名，leftNum为剩下的域名段，最高级为qq.com
func GetDomainHost(request *http.Request, leftNum int) string {
	if leftNum <= 2 {
		leftNum = 2
	}
	newHost := request.Header.Get("User-Host")
	if newHost == "" {
		newHost = request.Host
	}
	hostList := strings.Split(newHost, ".")
	newHostList := make([]string, 0)

	for _, one := range hostList {
		if one != "" {
			newHostList = append(newHostList, one)
		}
	}

	if len(newHostList) > leftNum {
		startNum := len(newHostList) - leftNum
		newHostList = newHostList[startNum:]
	}

	newHost = strings.Join(newHostList, ".")
	return newHost
}

// Ping 是否相通
func Ping(host string, port string, timeout time.Duration) error {
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Println("ERR:" + host + ">" + err.Error())
		return err
	}
	if conn != nil {
		_ = conn.Close()
		return nil
	}
	return fmt.Errorf("connect failed")
}

// GetLogId 生成唯一日志id
func GetLogId() string {
	logId := utils.NewUUID()
	randomStr := fmt.Sprintf("%s%s", logId, utils.GetRandomString(12))
	newLogId := utils.GetUUID(randomStr)
	logIdFront := newLogId[0:24]

	logIdEnd := conv.FormatFromUnixTime(conv.ShortTimeForm12)

	return strings.ToLower(fmt.Sprintf("%s%s", logIdFront, logIdEnd))
}
