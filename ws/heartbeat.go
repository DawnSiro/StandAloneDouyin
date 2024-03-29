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

var Heartbeats = make(map[*Client]struct{})

func (h *Hub) RunHeartbeatCheck() {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 对所有连接进行心跳检测
			for client := range Heartbeats {
				// 检查心跳是否超时
				if time.Since(client.LastHeartbeatTime) > TimeoutDuration {
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
