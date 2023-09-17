package model

import (
	"douyin/pkg/constant"
	"time"
)

type Video struct {
	ID            uint64    `json:"id"`
	PublishTime   time.Time `gorm:"not null" json:"publish_time"`
	AuthorID      uint64    `gorm:"not null" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(255);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(255);not null" json:"cover_url"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
	Title         string    `gorm:"type:varchar(63);not null" json:"title"`
}

func (n *Video) TableName() string {
	return constant.VideoTableName
}

type VideoData struct {
	// 没有ID的话，貌似第一个字段会被识别为主键ID
	VID               uint64 `gorm:"column:vid"`
	PlayURL           string
	CoverURL          string
	FavoriteCount     int64
	CommentCount      int64
	IsFavorite        bool
	Title             string
	UID               uint64
	Username          string
	FollowCount       int64
	FollowerCount     int64
	IsFollow          bool
	Avatar            string
	BackgroundImage   string
	Signature         string
	TotalFavorited    int64
	WorkCount         int64
	UserFavoriteCount int64
}
