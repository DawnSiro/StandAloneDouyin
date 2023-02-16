package db

import (
	"douyin/constant"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	UpdatedAt     time.Time `gorm:"column:update_time;not null;index:idx_update" `
	AuthorID      uint64    `gorm:"index:idx_authorid;not null"`
	PlayURL       string    `gorm:"type:varchar(255);not null"`
	CoverURL      string    `gorm:"type:varchar(255)"`
	FavoriteCount int64     `gorm:"default:0"`
	CommentCount  int64     `gorm:"default:0"`
	Title         string    `gorm:"type:varchar(50);not null"`
}

func (n *Video) TableName() string {
	return constant.VideoTableName
}

func CreateVideo(video *Video) error {
	return DB.Create(&video).Error
}

// MGetVideos multiple get list of videos info
func MGetVideos(maxVideoNum int, latestTime *int64) ([]*Video, error) {
	res := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}

	if err := DB.Where("update_time < ?", time.UnixMilli(*latestTime)).Limit(maxVideoNum).
		Order("update_time desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetVideosByAuthorID(userID uint64) ([]*Video, error) {
	res := make([]*Video, 0)
	err := DB.Find(&res, "author_id = ? ", userID).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectAuthorIDByVideoID(videoID int64) (uint64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
		},
	}

	result := DB.First(&video)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return video.AuthorID, nil
}

func UpdateFavoriteCount(videoID uint64, favoriteCount int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
		},
	}

	if err := DB.Model(&video).Update("favorite_count", favoriteCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseFavoriteCount increase 1
func IncreaseFavoriteCount(videoID uint64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
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
func ReduceFavoriteCount(videoID uint64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
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

func SelectFavoriteCountByVideoID(videoID int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
		},
	}

	result := DB.First(&video)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return video.FavoriteCount, nil
}

func SelectCommentCountByVideoID(videoID int64) (int64, error) {
	video := &Video{
		Model: gorm.Model{
			ID: uint(videoID),
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
