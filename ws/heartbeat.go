package ws

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"
	"time"
)

const (
	HeartbeatInterval = 5 * time.Second  // 心跳消息发送间隔
	TimeoutDuration   = 30 * time.Second // 连接超时时间
	//PongWait          = 10 * time.Second
)

func (h *Hub) GetConnectedClients() []*Client {
	clients := make([]*Client, 0)
	for _, client := range h.Clients {
		clients = append(clients, client)
	}
	return clients
}

var heartbeats = make(map[*Client]struct{})

// Todo: 心跳检测，超时重连
//func (c *Client) Reconnect() error {
//	// 创建新的 WebSocket 连接
//	newConn, _, err := websocket.DefaultDialer.Dial("ws://your-server-address", nil)
//	if err != nil {
//		return err // 返回连接错误
//	}
//
//	// 关闭旧连接
//	err = c.Conn.Close()
//	if err != nil {
//		return err // 返回关闭连接错误
//	}
//
//	// 更新连接
//	c.Conn = newConn
//	c.lastHeartbeatTime = time.Now()
//
//	return nil // 返回 nil 表示重新连接成功
//}

func (h *Hub) RunHeartbeatCheck() {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 对所有连接进行心跳检测
			for client := range heartbeats {
				// 检查心跳是否超时
				if time.Since(client.lastHeartbeatTime) > TimeoutDuration {
					hlog.Info("Heartbeat timeout, closing connection for client:", client.ID)
					h.Unregister <- client
					_ = client.Conn.Close()
				}

				err := client.Conn.WriteMessage(websocket.PingMessage, nil)

				if err != nil {
					hlog.Error("Error sending heartbeat ping:", err)
				}
			}
		}
	}
}
