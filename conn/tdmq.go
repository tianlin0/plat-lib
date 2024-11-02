// Package conn 连接参数
package conn

import (
	"fmt"
	"net"
	"strconv"
)

const _TDMQ_PROTOCOL = "pulsar"

type tdmConnect struct {
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

func (con *tdmConnect) getConnect() string {
	protocol := _TDMQ_PROTOCOL
	if con.Protocol != "" {
		protocol = con.Protocol
	}
	if con.Port == 0 {
		return fmt.Sprintf("%s://%s", protocol, con.Host)
	}
	tdmqURL := net.JoinHostPort(con.Host, strconv.Itoa(con.Port))
	return fmt.Sprintf("%s://%s", protocol, tdmqURL)
}
