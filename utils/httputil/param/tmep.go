package param

import (
	"net"
	"net/http"
)

const (
	// LocalHostIP localhost ip
	localHostIP = "127.0.0.1"
)

// GetInternalIPv4Address 获取内部IPv4地址
func GetInternalIPv4Address() string {
	addrStr, err := net.InterfaceAddrs()
	if err != nil {
		return localHostIP
	}
	for _, addr := range addrStr {
		ipaddr, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if ipaddr.IsLoopback() {
			continue
		}
		if ipaddr.To4() != nil {
			return ipaddr.String()
		}
	}
	return localHostIP
}

// ClientIP 得到客户端IP地址
func ClientIP(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

//// GetPathParam Returns a string parameter from request path or req.Attributes
//func GetPathParam(name string, req *http.Request) (param string) {
//	restfulReq := restful.NewRequest(req)
//	// Get parameter from request path
//	param = restfulReq.PathParameter(name)
//	if param != "" {
//		return param
//	}
//
//	// Get parameter from request attributes (set by intermediates)
//	attr := restfulReq.Attribute(name)
//	if attr != nil {
//		param, _ = attr.(string)
//	}
//	return
//}
//
//// GetIntParam 取得int参数
//func GetIntParam(req *http.Request, name string, def int) int {
//	restfulReq := restful.NewRequest(req)
//	num := def
//	if strNum := restfulReq.QueryParameter(name); len(strNum) > 0 {
//		num, _ = strconv.Atoi(strNum)
//	}
//	return num
//}
