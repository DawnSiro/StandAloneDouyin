package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func Messages(messages []*db.Message) []*api.Message {
	res := make([]*api.Message, 0)
	for i := 0; i < len(messages); i++ {
		res = append(res, Message(messages[i]))
	}
	return res
}

func Message(message *db.Message) *api.Message {
	createTime := message.SendTime.Unix()
	return &api.Message{
		ID:         int64(message.ID),
		ToUserID:   int64(message.ToUserID),
		FromUserID: int64(message.FromUserID),
		Content:    message.Content,
		CreateTime: &createTime,
	}
}
