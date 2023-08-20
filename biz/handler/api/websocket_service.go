package api

import (
	"context"
	"douyin/dal/db"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"strings"
)

type Client struct {
	ID       string
	ToUserID string
	Conn     *websocket.Conn
	Send     chan []byte
}

func (c *Client) readPump() {
	defer func() {
		MannaClient.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			hlog.Error("api.websocket_service.readPump.websocket_service_close err:", err.Error())
		}
	}()
	for {
		//c.Conn.PongHandler()
		SendMsg := new(SendMsg)

		err := c.Conn.ReadJSON(&SendMsg)
		if err != nil {
			hlog.Error("api.websocket_service.readPump.ReadJSON err:", err.Error())
			MannaClient.Unregister <- c
			_ = c.Conn.Close()
			break
		}

		if SendMsg.Type == 1 { // 发送消息
			MannaClient.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(SendMsg.Content), // 发送过来的消息
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
		case message, ok := <-c.Send: // 拿到Client的消息
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					hlog.Error("api.websocket_service.writePump.WriteMessage err:", err.Error())
					return
				}
				return
			}
			ReplyMsg := ReplyMsg{
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(ReplyMsg)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)

		case message, ok := <-c.Send: // 不在线逻辑
			if ok {
				uid, touid, err := ExtractNumbers(c.ToUserID)
				if err != nil {
					hlog.Error("api.websocket_service.writePump.ExtractNumbers err:", err.Error())
					return
				}
				isFriend := db.IsFriend(uid, touid)
				if !isFriend {
					errNo := errno.UserRequestParameterError
					errNo.ErrMsg = "不能给非好友发消息"
					//fmt.Printf("%d,%d\n", uid, touid)
					hlog.Error("api.websocket_service.writePump.IsFriend err:", errNo.Error())
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
// @router /douyin/message/ws/ [POST]
func ServeWs(ctx context.Context, c *app.RequestContext) {
	fromUserID := c.GetUint64(global.Config.JWTConfig.IdentityKey)
	hlog.Info("biz.handler.api.websocket_service.ServeWs GetFromUserID err:", fromUserID)
	toUid := c.Query("to_user_id")

	if strconv.FormatUint(fromUserID, 10) == toUid { // 不能给自己发消息
		hlog.Info("biz.handler.api.websocket_service.ServeWs FromUserID == toUid err:", fromUserID)
		return
	}

	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		client := &Client{
			ID:       CreateID(strconv.FormatUint(fromUserID, 10), toUid),
			ToUserID: CreateID(toUid, strconv.FormatUint(fromUserID, 10)),
			Conn:     conn,
			Send:     make(chan []byte),
		}

		MannaClient.Register <- client

		go client.writePump()
		client.readPump()
	})
	if err != nil {
		hlog.Error("biz.handler.api.websocket_service.ServeWs err:", err.Error())
	}
}
