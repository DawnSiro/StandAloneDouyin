package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
)

func SendMessage(fromUserID, toUserID uint64, content string) (*api.DouyinMessageActionResponse, error) {
	err := db.CreateMessage(fromUserID, toUserID, content)
	if err != nil {
		return nil, err
	}
	return &api.DouyinMessageActionResponse{StatusCode: 0}, nil
}

func GetMessageChat(userID, oppositeID uint64) (*api.DouyinMessageChatResponse, error) {
	messages, err := db.GetMessagesByUserID(userID, oppositeID)
	if err != nil {
		return nil, err
	}
	return &api.DouyinMessageChatResponse{
		StatusCode:  0,
		StatusMsg:   nil,
		MessageList: pack.Messages(messages),
	}, nil
}
