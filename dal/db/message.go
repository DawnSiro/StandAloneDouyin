package db

import (
	"douyin/dal/model"
	"douyin/dal/rdb"
	"douyin/pkg/global"
	"gorm.io/gorm"
	"time"

	"douyin/pkg/constant"
)

func CreateMessage(fromUserID, toUserID uint64, content string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		m := &model.Message{FromUserID: fromUserID, ToUserID: toUserID, Content: content, CreatedTime: time.Now()}
		err := global.DB.Create(m).Error
		if err != nil {
			return nil
		}
		return rdb.AddMessage(fromUserID, toUserID, m)
	})
}

func GetMessagesByUserIDAndPreMsgTime(userID, oppositeID uint64, preMsgTime int64) ([]*model.Message, error) {
	res := make([]*model.Message, 0)
	message := &model.Message{}
	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := global.DB.Raw("? UNION ? ORDER BY create_time ASC",
		global.DB.Where("to_user_id = ? AND from_user_id = ? AND `create_time` > ?",
			userID, oppositeID, time.UnixMilli(preMsgTime)).Model(message),
		global.DB.Where("to_user_id = ? AND from_user_id = ? AND `create_time` > ?",
			oppositeID, userID, time.UnixMilli(preMsgTime)).Model(message),
	).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetLatestMsg(userID uint64, oppositeID uint64) (*model.FriendMessageResp, error) {
	message := &model.Message{}
	// 使用 Union 来避免使用 or 导致不走索引的问题
	err := global.DB.Raw("? UNION ? ORDER BY create_time DESC LIMIT 1",
		global.DB.Where("to_user_id = ? AND from_user_id = ?", userID, oppositeID).Model(message),
		global.DB.Where("to_user_id = ? AND from_user_id = ?", oppositeID, userID).Model(message),
	).Scan(&message).Error
	if err != nil {
		return nil, err
	}

	switch message.ToUserID {
	case oppositeID:
		return &model.FriendMessageResp{
			Content: message.Content,
			MsgType: constant.SentMessage,
		}, nil
	default: // 默认发给自己
		return &model.FriendMessageResp{
			Content: message.Content,
			MsgType: constant.ReceivedMessage,
		}, nil
	}
}
