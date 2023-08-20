package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Messages(messages []*db.Message) []*api.Message {
	if messages == nil {
		hlog.Error("pack.message.Messages err:", errno.ServiceError)
		return nil
	}
	res := make([]*api.Message, 0)
	for i := 0; i < len(messages); i++ {
		res = append(res, Message(messages[i]))
	}
	return res
}

func Message(message *db.Message) *api.Message {
	if message == nil {
		hlog.Error("pack.message.Messages err:", errno.ServiceError)
		return nil
	}
	createTime := message.CreateTime.UnixMilli()
	return &api.Message{
		ID:         int64(message.ID),
		ToUserID:   int64(message.ToUserID),
		FromUserID: int64(message.FromUserID),
		Content:    message.Content,
		CreateTime: &createTime,
	}
}
