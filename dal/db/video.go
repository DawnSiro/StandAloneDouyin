package db

import (
	"time"

	"douyin/pkg/constant"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
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

func CreateVideo(video *Video) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		//从这里开始，应该使用 tx 而不是 db（tx 是 Transaction 的简写）
		u := &User{ID: video.AuthorID}
		err := tx.Select("work_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("work_count", u.WorkCount+1).Error
		if err != nil {
			return err
		}
		err = tx.Create(video).Error
		if err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

// MGetVideos multiple get list of videos info
func MGetVideos(maxVideoNum int, latestTime *int64) ([]*Video, error) {
	res := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		// 这里和文档里说得不一样，实际客户端传的是毫秒
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}

	// TODO 设计简单的推荐算法，比如关注的 UP 发了视频，会优先推送
	hlog.Info(*latestTime)
	if err := DB.Where("publish_time < ?", time.UnixMilli(*latestTime)).Limit(maxVideoNum).
		Order("publish_time desc").Find(&res).Error; err != nil {
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

func SelectAuthorIDByVideoID(videoID uint64) (uint64, error) {
	video := &Video{
		ID: videoID,
	}

	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.AuthorID, nil
}

func UpdateVideoFavoriteCount(videoID uint64, favoriteCount uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	if err := DB.Model(&video).Update("favorite_count", favoriteCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseVideoFavoriteCount increase 1
func IncreaseVideoFavoriteCount(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}
	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := DB.Model(&video).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

// DecreaseVideoFavoriteCount decrease 1
func DecreaseVideoFavoriteCount(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}
	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := DB.Model(&video).Update("favorite_count", video.FavoriteCount-1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

func UpdateCommentCount(videoID uint64, commentCount uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	if err := DB.Model(&video).Update("comment_count", commentCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseCommentCount increase 1
func IncreaseCommentCount(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}

	if err := DB.Model(&video).Update("comment_count", video.CommentCount+1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

// DecreaseCommentCount decrease  1
func DecreaseCommentCount(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}
	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := DB.Model(&video).Update("comment_count", video.CommentCount-1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

func SelectVideoFavoriteCountByVideoID(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	err := DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

func SelectCommentCountByVideoID(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	err := DB.First(&video).Error
	if err != nil {
		return 0, err
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
