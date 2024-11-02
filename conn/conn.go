// Package conn 连接参数
package conn

import (
	"net"
	"strconv"
)

const (
	DriverMysql = "mysql"
	DriverRedis = "redis"
	DriverTdmq  = "tdmq"
)

type connInterface interface {
	getConnect() string //获取连接字符串
}

var _, _, _ connInterface = new(mysqlConnect), new(redisConnect), new(tdmConnect)

// Connect 数据连接对象
type Connect struct {
	Driver   string                 `json:"driver,omitempty"`
	Protocol string                 `json:"protocol,omitempty"`
	Host     string                 `json:"host,omitempty"`
	Port     string                 `json:"port,omitempty"`
	Username string                 `json:"username,omitempty"`
	Password string                 `json:"password,omitempty"`
	Database string                 `json:"database,omitempty"`
	Extend   map[string]interface{} `json:"extend,omitempty"`
}

// ConnFunc 数据连接参数
type connOpt func(*Connect)

// NewOption 新增
func NewOption() connOpt {
	return func(*Connect) {}
}

// DialDriver 连接类型
func (c connOpt) DialDriver(driver string) connOpt {
	return func(do *Connect) {
		c(do)
		do.Driver = driver
	}
}

// DialProtocol 连接协议
func (c connOpt) DialProtocol(protocol string) connOpt {
	return func(do *Connect) {
		c(do)
		do.Protocol = protocol
	}
}

// DialHostPort 连接ip和端口号
func (c connOpt) DialHostPort(host string, port string) connOpt {
	return func(do *Connect) {
		c(do)
		do.Host = host
		do.Port = port
	}
}

// DialDatabase 连接库
func (c connOpt) DialDatabase(db string) connOpt {
	return func(do *Connect) {
		c(do)
		do.Database = db
	}
}

// DialUserNamePassword 连接用户名和密码
func (c connOpt) DialUserNamePassword(username, password string) connOpt {
	return func(do *Connect) {
		c(do)
		do.Username = username
		do.Password = password
	}
}

// DialExtend 扩展函数
func (c connOpt) DialExtend(ext map[string]interface{}) connOpt {
	return func(do *Connect) {
		c(do)

		if do.Extend == nil {
			do.Extend = make(map[string]interface{})
		}
		if ext == nil {
			return
		}

		for k, v := range ext {
			do.Extend[k] = v
		}
	}
}

// GetConnString 获取连接字符串，不同的driver，返回的不同
func (con *Connect) GetConnString(cp ...connOpt) string {
	for _, one := range cp {
		one(con)
	}
	return getConnString(con)
}

// getConnString 获取连接字符串，不同的driver，返回的不同
func getConnString(con *Connect) string {
	if con == nil {
		return ""
	}

	if con.Driver == DriverMysql {
		myConn := &mysqlConnect{
			Host:     con.Host,
			Username: con.Username,
			Password: con.Password,
			Database: con.Database,
		}
		if con.Port != "" {
			p, err := strconv.Atoi(con.Port)
			if err == nil {
				myConn.Port = p
			}
		}
		if con.Extend != nil {
			if char, ok := con.Extend["charset"]; ok {
				myConn.Charset = char.(string)
			}
		}
		return myConn.getConnect()
	}

	if con.Driver == DriverRedis {
		redisConn := &redisConnect{
			Host:     con.Host,
			Username: con.Username,
			Password: con.Password,
		}
		if con.Port != "" {
			p, err := strconv.Atoi(con.Port)
			if err == nil {
				redisConn.Port = p
			}
		}
		if con.Database != "" {
			p, err := strconv.Atoi(con.Database)
			if err == nil {
				redisConn.Database = p
			}
		}
		return redisConn.getConnect()
	}

	if con.Driver == DriverTdmq {
		tdmConn := &tdmConnect{
			Protocol: con.Protocol,
			Host:     con.Host,
		}
		if con.Port != "" {
			p, err := strconv.Atoi(con.Port)
			if err == nil {
				tdmConn.Port = p
			}
		}
		return tdmConn.getConnect()
	}

	//默认的地址
	if con.Host != "" {
		if con.Port != "" {
			return net.JoinHostPort(con.Host, con.Port)
		}
		return con.Host
	}

	return ""
}
