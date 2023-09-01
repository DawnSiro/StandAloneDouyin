package ws

import (
	"context"
	"douyin/biz/service"
	"douyin/dal/db"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
)

type Client struct {
	ID                string
	ToUserID          string
	Conn              *websocket.Conn
	Send              chan []byte
	LastHeartbeatTime time.Time
	ConnMutex         sync.Mutex // 互斥锁用来保护连接操作
}

func (c *Client) readPump() {
	defer func() {
		MannaClient.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			hlog.Error("api.ws.websocket_service.readPump.websocket_service_close err:", err.Error())
		}
	}()
	for {
		//c.Conn.PongHandler()
		SendMsg := new(SendMsg)

		err := c.Conn.ReadJSON(&SendMsg)
		if err != nil {
			hlog.Error("api.ws.websocket_service.readPump.ReadJSON err:", err.Error())
			MannaClient.Unregister <- c
			_ = c.Conn.Close()
			break
		}

		if SendMsg.Type == 1 { // 发送消息
			MannaClient.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(SendMsg.Content), // 发送过来的消息
			}
		} else if SendMsg.Type == 2 { // 获取历史消息
			uid, touid, err := ExtractNumbers(c.ToUserID)
			if err != nil {
				hlog.Error("ws.hub.readPump.ExtractNumbers err:", err)
			}

			preMsgTime := int64(0)
			resp, err := service.GetMessageChat(uid, touid, preMsgTime)

			message := resp.MessageList
			messageCount := len(message)
			if messageCount > 100 {
				// 计算起始索引，以获取最后的100条消息
				startIndex := messageCount - 100
				for _, message := range message[startIndex:] {
					ReplyMsg := ReplyMsg{
						Content: message.Content,
					}
					//msg, _ := json.Marshal(ReplyMsg)
					_ = c.Conn.WriteJSON(ReplyMsg)
				}
			} else {
				for _, message := range message {
					ReplyMsg := ReplyMsg{
						Content: message.Content,
					}
					//msg, _ := json.Marshal(ReplyMsg)
					_ = c.Conn.WriteJSON(ReplyMsg)
				}
			}
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send: // 好友在线
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					hlog.Error("api.ws.websocket_service.writePump.WriteMessage err:", err.Error())
					return
				}
				return
			}
			ReplyMsg := ReplyMsg{
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(ReplyMsg)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)

		case message, ok := <-c.Send: // 好友不在线
			if ok {
				uid, touid, err := ExtractNumbers(c.ToUserID)
				if err != nil {
					hlog.Error("api.ws.websocket_service.writePump.ExtractNumbers err:", err.Error())
					return
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend {
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "Cannot send messages to non-friends"
					hlog.Error("api.ws.websocket_service.writePump.IsFriend err:", errNo.Error())
				}
				ReplyMsg := ReplyMsg{
					Content: fmt.Sprintf("%s", string(message)),
				}
				msg, _ := json.Marshal(ReplyMsg)
				_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

func CreateID(uid, toUid string) string {
	var builder strings.Builder
	builder.WriteString(uid)
	builder.WriteString("->")
	builder.WriteString(toUid)
	return builder.String()
}

// ServeWs .
// @router /douyin/message/ws/ [WebSocket]
func ServeWs(ctx context.Context, c *app.RequestContext) {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fromUserID := c.GetUint64(global.Config.JWTConfig.IdentityKey)
	hlog.Info("biz.handler.api.ws.websocket_service.ServeWs GetFromUserID:", fromUserID)
	toUid := c.Query("to_user_id")

	if strconv.FormatUint(fromUserID, 10) == toUid { // 不能给自己发消息
		hlog.Info("biz.handler.api.ws.websocket_service.ServeWs FromUserID == toUid err:", fromUserID)
		return
	}

	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		client := &Client{
			ID:                CreateID(strconv.FormatUint(fromUserID, 10), toUid),
			ToUserID:          CreateID(toUid, strconv.FormatUint(fromUserID, 10)),
			Conn:              conn,
			Send:              make(chan []byte),
			LastHeartbeatTime: time.Now(),
		}
		Heartbeats[client] = struct{}{}
		defer delete(Heartbeats, client)

		client.Conn.SetPongHandler(func(string) error {
			client.LastHeartbeatTime = time.Now()
			hlog.Info("Update client lastHeartbeatTime time")
			return nil
		})

		MannaClient.Register <- client

		go client.writePump()
		client.readPump()
	})
	if err != nil {
		hlog.Error("biz.handler.api.ws.websocket_service.ServeWs err:", err.Error())
	}
}
