package db

import (
	"douyin/constant"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ToUserId   uint64 `json:"to_user_id"`
	FromUserId uint64 `json:"from_user_id"`
	Content    string `json:"content"`
}

func (n *Message) TableName() string {
	return constant.MessageTableName
}

type FriendMessageResp struct {
	Content string
	MsgType uint64
}

func GetLatestMsg(userID uint64, toUserID uint64) (FriendMessageResp, error) {
	friendMessageResp := &FriendMessageResp{}
	message1 := &Message{}
	message2 := &Message{}

	//toUserId 当前用户id
	//fromUserId 发送信息用户
	//TODO: 开会需要理清楚

	con1 := DB.Where("to_user_id = ? and from_user_id = ?", userID, toUserID).First(&message1)
	con2 := DB.Where("to_user_id = ? and from_user_id = ?", toUserID, userID).First(&message2)

	if con1.RowsAffected != 0 && con2.RowsAffected != 0 {
		if message1.ID > message2.ID {
			friendMessageResp.MsgType = 1
			friendMessageResp.Content = message1.Content
			return *friendMessageResp, nil
		} else {
			friendMessageResp.MsgType = 0
			friendMessageResp.Content = message2.Content
			return *friendMessageResp, nil
		}
	} else if con1.RowsAffected != 0 {
		friendMessageResp.MsgType = 1
		friendMessageResp.Content = message1.Content
		return *friendMessageResp, nil
	} else if con2.RowsAffected != 0 {
		friendMessageResp.MsgType = 0
		friendMessageResp.Content = message2.Content
		return *friendMessageResp, nil
	}
	return FriendMessageResp{
		MsgType: 3,
	}, nil

}
