package model

import (
	"time"

	"douyin/pkg/constant"
)

// 统一存放数据模型结构体，避免 mysql 和 redis 之间出现循环引用问题

type Comment struct {
	ID          uint64    `json:"id"`
	IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	VideoID     uint64    `gorm:"not null" json:"video_id"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	Content     string    `gorm:"type:varchar(255);not null" json:"content"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

func (n *Comment) TableName() string {
	return constant.CommentTableName
}

type CommentData struct {
	CID            uint64 `gorm:"column:cid"`
	Content        string
	CreatedTime    time.Time
	UID            uint64
	Username       string
	FollowingCount uint64
	FollowerCount  uint64
	Avatar         string
	IsFollow       bool
}

type CommentRedisData struct {
	CID         uint64 `gorm:"column:cid"`
	Content     string
	CreatedTime time.Time
	UID         uint64
	Username    string
	Avatar      string
}
