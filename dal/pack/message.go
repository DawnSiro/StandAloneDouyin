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
	res := new(api.Message)
	res.ID = int64(message.ID)
	res.ToUserID = int64(message.ToUserID)
	res.FromUserID = int64(message.FromUserID)
	res.Content = message.Content
	createTime := message.CreatedAt.Unix()
	res.CreateTime = &createTime
	return res
}
