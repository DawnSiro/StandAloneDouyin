// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package api

import (
	"douyin/dal/db"
	"douyin/pkg/errno"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"strings"
)

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	Content string `json:"content"`
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

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
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
	for {
		hlog.Info("-----监听管道通信-----")
		select {
		case client := <-MannaClient.Register:
			MannaClient.Clients[client.ID] = client

			//fmt.Printf("有新连接: %v\n", client.ID) // 测试用
			ReplyMsg := ReplyMsg{
				Content: "已经连接到服务器了",
			}
			msg, _ := json.Marshal(ReplyMsg)
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)

			if err != nil {
				hlog.Error("biz.handler.api.hub.Run.WriteMessage err:", err.Error())
			}
		case client := <-h.Unregister:
			//fmt.Printf("连接失败%s,", client.ID)
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
			flag := false                       // 默认对方是不在线的 false表示不在线，ture为在线（用来标记消息是否已读）
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
				replyMsg := &ReplyMsg{
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Conn.WriteMessage(websocket.TextMessage, msg) // 对消息进行广播

				uid, touid, err := ExtractNumbers(broadcast.Client.ToUserID)
				if err != nil {
					//fmt.Println("Error:", err)
					hlog.Error("biz.handler.api.hub.ExtractNumbers err:", err)
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend { // 是好友将消息插入数据库，不是就退出
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "不能给非好友发消息"
					hlog.Error("biz.handler.api.hub.IsFriend err:", errNo.Error())

				} else {
					//fmt.Println(msg)
					ReplyMsg := ReplyMsg{
						Content: fmt.Sprintf("%s", string(message)),
					}
					msg, _ := json.Marshal(ReplyMsg)
					_ = broadcast.Client.Conn.WriteMessage(websocket.TextMessage, msg)

					err = db.CreateMessage(uid, touid, string(msg[12:len(msg)-2])) // 将消息放到数据库
					if err != nil {
						hlog.Error("biz.handler.api.hub.CreateMessage err:", err.Error())
					}
				}
			} else { // 好友不在线
				uid, touid, err := ExtractNumbers(broadcast.Client.ToUserID)
				if err != nil {
					//fmt.Println("Error:", err)
					hlog.Error("biz.handler.api.hub.ExtractNumbers err:", err)
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend { // 是好友将消息插入数据库，不是就退出
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "不能给非好友发消息"
					hlog.Error("biz.handler.api.hub.IsFriend err:", errNo.Error())

				} else {
					//fmt.Println(msg)
					ReplyMsg := ReplyMsg{
						Content: fmt.Sprintf("%s", string(message)),
					}
					msg, _ := json.Marshal(ReplyMsg)
					_ = broadcast.Client.Conn.WriteMessage(websocket.TextMessage, msg)

					err = db.CreateMessage(uid, touid, string(msg[12:len(msg)-2])) // 将消息放到数据库
					if err != nil {
						hlog.Error("biz.handler.api.hub.CreateMessage err:", err.Error())
					}
				}
			}
		}
	}
}
