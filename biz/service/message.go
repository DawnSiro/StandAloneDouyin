package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func SendMessage(fromUserID, toUserID uint64, content string) (*api.DouyinMessageActionResponse, error) {
	err := db.CreateMessage(fromUserID, toUserID, content)
	if err != nil {
		hlog.Error("service.message.SendMessage err:", err.Error())
		return nil, err
	}
	return &api.DouyinMessageActionResponse{StatusCode: 0}, nil
}

func GetMessageChat(userID, oppositeID uint64, preMsgTime int64) (*api.DouyinMessageChatResponse, error) {
	if userID == oppositeID {
		return nil, errno.UserRequestParameterError
	}
	messages, err := db.GetMessagesByUserIDAndPreMsgTime(userID, oppositeID, preMsgTime)
	if err != nil {
		hlog.Error("service.message.GetMessageChat err:", err.Error())
		return nil, err
	}
	return &api.DouyinMessageChatResponse{
		StatusCode:  0,
		MessageList: pack.Messages(messages),
	}, nil
}
