package model

import (
	"douyin/pkg/constant"
	"time"
)

type Message struct {
	ID          uint64    `json:"id"`
	ToUserID    uint64    `gorm:"not null" json:"to_user_id"`
	FromUserID  uint64    `gorm:"not null" json:"from_user_id"`
	Content     string    `gorm:"type:varchar(255);not null" json:"content"`
	CreatedTime time.Time `gorm:"not null" json:"created_time" `
}

func (n *Message) TableName() string {
	return constant.MessageTableName
}

type FriendMessageResp struct {
	Content string
	MsgType uint8
}
