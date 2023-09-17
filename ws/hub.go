// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package ws

import (
	"douyin/dal/db"
	"douyin/pkg/errno"
	"douyin/pkg/pulsar"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
)

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	Content string `json:"content"`
	//CreatedTime *int64 `json:"create_time" `
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var MannaClient = newHub()

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Reply:      make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func ExtractNumbers(s string) (uint64, uint64, error) {
	parts := strings.Split(s, "->")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Invalid format")
	}

	firstNum, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	secondNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return uint64(firstNum), uint64(secondNum), nil
}

func (h *Hub) Run() {
	go MannaClient.RunHeartbeatCheck() // 打开心跳检测

	for {
		hlog.Info("Monitor pipe communication")

		select {
		case client := <-MannaClient.Register:
			MannaClient.Clients[client.ID] = client

			ReplyMsg := ReplyMsg{
				Content: "Already connected to the websocket server",
			}
			msg, _ := json.Marshal(ReplyMsg)
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)

			if err != nil {
				hlog.Error("biz.handler.api.ws.hub.Run.WriteMessage err:", err.Error())
			}
		case client := <-h.Unregister:
			if _, ok := MannaClient.Clients[client.ID]; ok {
				ReplyMsg := &ReplyMsg{
					Content: "连接中断",
				}

				msg, _ := json.Marshal(ReplyMsg)
				_ = client.Conn.WriteMessage(websocket.TextMessage, msg)

				close(client.Send)
				delete(MannaClient.Clients, client.ID)
			}
		case broadcast := <-MannaClient.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.ToUserID // 2->1
			flag := false                       // 默认对方是不在线的 false 表示不在线，ture 为在线（用来标记消息是否已读）

			for id, conn := range MannaClient.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(MannaClient.Clients, conn.ID)
				}
			}

			if flag {
				uid, touid, err := ExtractNumbers(broadcast.Client.ToUserID)
				if err != nil {
					hlog.Error("biz.handler.api.ws.hub.ExtractNumbers err:", err)
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend { // 是好友将消息插入数据库，不是就退出
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "Cannot send messages to non-friends"
					hlog.Error("biz.handler.api.ws.hub.IsFriend err:", errNo.Error())
				} else {
					broadcast.Client.LastHeartbeatTime = time.Now() // 更新心跳时间
					hlog.Info("Update client lastHeartbeatTime time")

					// 将消息发布到消息队列
					err = pulsar.GetMessageMQInstance().CreateMessage(uid, touid, string(message))
					if err != nil {
						hlog.Error("biz.handler.api.ws.hub.CreateMessage publish message err: ", err.Error())
					}
				}
			} else { // 好友不在线
				uid, touid, err := ExtractNumbers(broadcast.Client.ToUserID)
				if err != nil {
					hlog.Error("biz.handler.api.ws.hub.ExtractNumbers err:", err)
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend { // 是好友将消息插入数据库，不是就退出
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "Cannot send messages to non-friends"
					hlog.Error("biz.handler.api.ws.hub.IsFriend err:", errNo.Error())

				} else {
					broadcast.Client.LastHeartbeatTime = time.Now() // 更新心跳时间
					hlog.Info("Update client lastHeartbeatTime time")

					// 将消息发布到消息队列
					err = pulsar.GetMessageMQInstance().CreateMessage(uid, touid, string(message))
					if err != nil {
						hlog.Error("biz.handler.api.ws.hub.CreateMessage publish message err: ", err.Error())
					}
				}
			}
		}
	}
}
