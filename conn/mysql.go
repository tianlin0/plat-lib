// Package conn 连接参数
package conn

import (
	"fmt"
	"net/url"
)

const _MYSQL_CHARSET = "utf8"

// mysqlConnect 数据连接对象
type mysqlConnect struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Charset  string `json:"charset,omitempty"`
}

func getMysqlCharset(con *mysqlConnect) string {
	if con.Charset != "" {
		return con.Charset
	}
	return _MYSQL_CHARSET
}

func (con *mysqlConnect) getConnect() string {
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		con.Username,
		url.QueryEscape(con.Password),
		con.Host,
		con.Port,
		con.Database,
		getMysqlCharset(con))
	return conStr
}
