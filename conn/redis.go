// Package conn 连接参数
package conn

import (
	"fmt"
	"net/url"
)

type redisConnect struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database int    `json:"database,omitempty"`
}

func (con *redisConnect) getConnect() string {
	conStr := fmt.Sprintf("redis://%s:%s@%s:%d/%d",
		con.Username,
		url.QueryEscape(con.Password),
		con.Host,
		con.Port,
		con.Database)
	return conStr
}
