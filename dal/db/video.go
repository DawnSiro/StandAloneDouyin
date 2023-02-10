package db

import (
	"douyin/constant"
	"errors"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	AuthorId      uint64 `json:"author_id"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	Title         string `json:"title"`
}

func (n *Video) TableName() string {
	return constant.VideoTableName
}

func SelectAuthorIdByVideoId(videoId int64) (uint64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}

	result := DB.First(&video)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return video.AuthorId, nil
}
