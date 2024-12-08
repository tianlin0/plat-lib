package websocket

import (
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// PostMessage 启动客户端
func PostMessage(ctx context.Context, wsUrl string, msg string) (string, error) {
	// 连接到 WebSocket 服务器 "ws://localhost:8080"
	conn, _, _, err := ws.Dial(ctx, wsUrl)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 发送消息到服务器
	if err = wsutil.WriteClientMessage(conn, ws.OpText, []byte(msg)); err != nil {
		return "", err
	}

	// 接收服务器的回复
	retMsg, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		return "", err
	}

	return string(retMsg), nil
}
