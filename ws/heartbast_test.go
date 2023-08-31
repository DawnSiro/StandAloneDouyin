package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"testing"
)

func TestRunHeartbeatCheck(t *testing.T) {
	url := "ws://127.0.0.1:8080/douyin/message/ws/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjkyODEyNzE0LCJuYmYiOjE2OTI3Njk1MTQsImlhdCI6MTY5Mjc2OTUxNH0.JMEbBafvt8olZRWj4U0vG-79nwfHIpahBKmxwt2S6Oc&to_user_id=4"

	var conn *websocket.Conn
	var err error

	connectWebSocket := func() {
		conn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Println("Error connecting to WebSocket server:", err)
			return
		}
	}

	connectWebSocket() // Initial connection

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	var writeMutex sync.Mutex

	//// 启动一个 goroutine 模拟心跳, 客户端发送验证
	//go func() {
	//	for {
	//		writeMutex.Lock()
	//		err := conn.WriteMessage(websocket.PingMessage, nil)
	//		writeMutex.Unlock()
	//
	//		if err != nil {
	//			fmt.Println("Error sending ping:", err)
	//			for retry := 0; retry < 3; retry++ {
	//				fmt.Println("Retrying connection...")
	//				connectWebSocket()
	//				time.Sleep(5 * time.Second)
	//				writeMutex.Lock()
	//				err := conn.WriteMessage(websocket.PingMessage, nil)
	//				writeMutex.Unlock()
	//				if err == nil {
	//					break
	//				}
	//			}
	//		} else {
	//			fmt.Println("Websocket Keep-Alive")
	//		}
	//
	//		time.Sleep(5 * time.Second)
	//	}
	//}()

	// 在连接正常的情况下，可以模拟发送和接收消息
	go func() {
		for {
			sendMsg := SendMsg{
				Type:    1, // 设置消息类型
				Content: "Hello, WebSocket Server!",
			}

			writeMutex.Lock()
			err := conn.WriteJSON(sendMsg)
			writeMutex.Unlock()
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}

			// 接收消息
			var receivedMsg SendMsg
			err = conn.ReadJSON(&receivedMsg)
			if err != nil {
				fmt.Println("Error reading message:", err)
				break
			}
			fmt.Println("Received message:", receivedMsg.Content)
		}
	}()

	// 保持主程序运行，不关闭连接
	select {}
}
