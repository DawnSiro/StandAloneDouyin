package db

import (
	"douyin/constants"
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
	return constants.VideoTableName
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

func UpdateFavoriteCount(videoId uint64, favoriteCount int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}

	if err := DB.Model(&video).Update("favorite_count", favoriteCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseFavoriteCount increase 1
func IncreaseFavoriteCount(videoId uint64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}
	if err := DB.Find(&video).Error; err != nil {
		return 0, err
	}
	if err := DB.Model(&video).Update("comment_count", video.CommentCount+1).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// ReduceFavoriteCount reduce 1
func ReduceFavoriteCount(videoId uint64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}
	if err := DB.Find(&video).Error; err != nil {
		return 0, err
	}
	if err := DB.Model(&video).Update("comment_count", video.CommentCount-1).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

func SelectFavoriteCountByVideoId(videoId int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}

	result := DB.First(&video)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return video.FavoriteCount, nil
}

func SelectCommentCountByVideoId(videoId int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoId),
		},
	}

	result := DB.First(&video)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return video.CommentCount, nil
}

func SelectVideoList() ([]*Video, error) {
	videoList := new([]*Video)

	if err := DB.Find(&videoList).Error; err != nil {
		return nil, err
	}
	return *videoList, nil
}
