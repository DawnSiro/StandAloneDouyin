package db

import (
	"douyin/constant"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ToUserID   uint64 `json:"to_user_id"`
	FromUserID uint64 `json:"from_user_id"`
	Content    string `json:"content"`
}

func (n *Message) TableName() string {
	return constant.MessageTableName
}

type FriendMessageResp struct {
	Content string
	MsgType uint64
}

func CreateMessage(fromUserID, toUserID uint64, content string) error {
	return DB.Create(&Message{FromUserID: fromUserID, ToUserID: toUserID, Content: content}).Error
}

func GetMessagesByUserID(userID, oppositeID uint64) ([]*Message, error) {
	res := make([]*Message, 0)

	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := DB.Raw("? UNION ? ORDER BY id DESC",
		DB.Where("to_user_id = ? AND from_user_id = ?", userID, oppositeID).Model(&Message{}),
		DB.Where("to_user_id = ? AND from_user_id = ?", oppositeID, userID).Model(&Message{}),
	).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetLatestMsg(userID uint64, toUserID uint64) (*FriendMessageResp, error) {
	message := &Message{}

	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := DB.Raw("? UNION ? ORDER BY id DESC LIMIT 1",
		DB.Where("to_user_id = ? AND from_user_id = ?", userID, toUserID).Model(&Message{}),
		DB.Where("to_user_id = ? AND from_user_id = ?", toUserID, userID).Model(&Message{}),
	).Scan(&message).Error
	if err != nil {
		return nil, err
	}

	switch message.ToUserID {
	case toUserID:
		return &FriendMessageResp{
			Content: message.Content,
			MsgType: 1,
		}, nil
	default: // 默认发给自己
		return &FriendMessageResp{
			Content: message.Content,
			MsgType: 0,
		}, nil
	}

}
