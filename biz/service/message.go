package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"douyin/dal/pack"
	"errors"
)

func SendMessage(fromUserID, toUserID uint64, actionType int32, content string) (*api.DouyinMessageActionResponse, error) {
	if actionType == constant.SendMessageAction {
		err := db.CreateMessage(fromUserID, toUserID, content)
		if err != nil {
			return nil, err
		}
		return &api.DouyinMessageActionResponse{StatusCode: 0}, nil
	}
	return nil, errors.New("action type error")
}

func GetMessageChat(userID, oppositeID uint64) (*api.DouyinMessageChatResponse, error) {
	res := new(api.DouyinMessageChatResponse)

	messages, err := db.GetMessagesByUserID(userID, oppositeID)
	if err != nil {
		return nil, err
	}
	messageList := pack.Messages(messages)
	res.MessageList = messageList
	return res, nil
}
