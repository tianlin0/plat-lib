package websocket

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
)

// handleConnection 处理客户端请求
func handleConnection(conn net.Conn) error {
	defer conn.Close()

	// 循环接收消息
	for {
		// 接收客户端消息
		msg, op, err := wsutil.ReadServerData(conn)
		if err != nil {
			log.Println("读取消息出错:", err)
			return err
		}

		// 打印收到的消息
		fmt.Printf("收到消息: %s\n", string(msg))

		// 生成一个随机数并发送回客户端
		randomNumber := fmt.Sprintf("随机数：%d", 42) // 这里可以随机生成数字
		if err := wsutil.WriteServerMessage(conn, op, []byte(randomNumber)); err != nil {
			log.Println("发送消息出错:", err)
			return err
		}
	}
	return nil
}

// StartService 设置 WebSocket 服务端
func StartService(r *http.Request, w http.ResponseWriter) error {
	// 升级 HTTP 请求为 WebSocket 连接
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return err
	}

	// 处理连接
	return handleConnection(conn)
}
