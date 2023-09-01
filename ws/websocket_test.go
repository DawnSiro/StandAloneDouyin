package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func TestServeWs(t *testing.T) {
	// 创建WebSocket连接
	//u := url.URL{Scheme: "ws", Host: "localhost:30000", Path: "/ws"}
	u := "ws://127.0.0.1:30000/douyin/message/ws/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NCwiZXhwIjoxNjkzNjIzMDg0LCJuYmYiOjE2OTM1Nzk4ODQsImlhdCI6MTY5MzU3OTg4NH0.3zPhg8nRfPqVYQRzGSoQPX9QTQGlIbL404XUTTXqm-4&to_user_id=2"
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Error connecting to WebSocket server: %v", err)
	}

	// 发送消息
	go func() {
		for i := 0; i < 1000; i++ {
			sendMsg := SendMsg{
				Type:    1,
				Content: fmt.Sprintf("Hello, WebSocket Server! (%d)", i),
			}
			time.Sleep(time.Millisecond) // 控制发送速率，避免过于快速

			err := conn.WriteJSON(sendMsg)
			if err != nil {
				t.Fatalf("Error sending message: %v", err)
			}
		}
	}()

	select {}
}
