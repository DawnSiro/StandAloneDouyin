package model

import (
	"douyin/pkg/constant"
	"time"
)

type Relation struct {
	ID          uint64    `json:"id"`
	IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	ToUserID    uint64    `gorm:"not null" json:"to_user_id"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

func (n *Relation) TableName() string {
	return constant.RelationTableName
}

// FollowUserData Redis ZSet 中使用 CreatedTime 做为 score
type FollowUserData struct {
	UID         uint64 `gorm:"column:uid"`
	Username    string
	Avatar      string
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

// FollowUserRedisData Redis ZSet 中需要序列化存储的 Member
type FollowUserRedisData struct {
	UID      uint64 `gorm:"column:uid"`
	Username string
	Avatar   string
}

type FanUserData struct {
	UID         uint64 `gorm:"column:uid"`
	Username    string
	Avatar      string
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

type FanUserRedisData struct {
	UID      uint64 `gorm:"column:uid"`
	Username string
	Avatar   string
}

type RelationUserData struct {
	UID            uint64 `gorm:"column:uid"`
	Username       string
	FollowingCount uint64
	FollowerCount  uint64
	Avatar         string
	IsFollow       bool
}

type FriendUserData struct {
	ID       uint64 // 用户id
	Name     string // 用户名称
	IsFollow bool   // true-已关注，false-未关注
	Avatar   string // 用户头像Url
	Message  string // 和该好友的最新聊天消息
	// message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	MsgType int8
}
