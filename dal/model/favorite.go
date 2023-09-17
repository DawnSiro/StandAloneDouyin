package model

import (
	"douyin/pkg/constant"
	"time"
)

type UserFavoriteVideo struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	VideoID     uint64    `gorm:"not null" json:"video_id"`
	IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

func (n *UserFavoriteVideo) TableName() string {
	return constant.UserFavoriteVideosTableName
}
