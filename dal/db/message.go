package db

import (
	"time"

	"douyin/pkg/constant"
)

type Message struct {
	ID         uint64    `json:"id"`
	ToUserID   uint64    `gorm:"not null" json:"to_user_id"`
	FromUserID uint64    `gorm:"not null" json:"from_user_id"`
	Content    string    `gorm:"type:varchar(255);not null" json:"content"`
	CreateTime time.Time `gorm:"not null" json:"create_time" `
}

func (n *Message) TableName() string {
	return constant.MessageTableName
}

type FriendMessageResp struct {
	Content string
	MsgType uint8
}

func CreateMessage(fromUserID, toUserID uint64, content string) error {
	return DB.Create(&Message{FromUserID: fromUserID, ToUserID: toUserID, Content: content, CreateTime: time.Now()}).Error
}

func GetMessagesByUserIDAndPreMsgTime(userID, oppositeID uint64, preMsgTime int64) ([]*Message, error) {
	res := make([]*Message, 0)
	message := &Message{}
	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := DB.Raw("? UNION ? ORDER BY create_time ASC",
		DB.Where("to_user_id = ? AND from_user_id = ? AND `create_time` > ?",
			userID, oppositeID, time.UnixMilli(preMsgTime+100)).Model(message),
		DB.Where("to_user_id = ? AND from_user_id = ? AND `create_time` > ?",
			oppositeID, userID, time.UnixMilli(preMsgTime+100)).Model(message),
	).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetLatestMsg(userID uint64, oppositeID uint64) (*FriendMessageResp, error) {
	message := &Message{}
	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := DB.Raw("? UNION ? ORDER BY create_time DESC LIMIT 1",
		DB.Where("to_user_id = ? AND from_user_id = ?", userID, oppositeID).Model(message),
		DB.Where("to_user_id = ? AND from_user_id = ?", oppositeID, userID).Model(message),
	).Scan(&message).Error
	if err != nil {
		return nil, err
	}

	switch message.ToUserID {
	case oppositeID:
		return &FriendMessageResp{
			Content: message.Content,
			MsgType: constant.SentMessage,
		}, nil
	default: // 默认发给自己
		return &FriendMessageResp{
			Content: message.Content,
			MsgType: constant.ReceivedMessage,
		}, nil
	}
}
