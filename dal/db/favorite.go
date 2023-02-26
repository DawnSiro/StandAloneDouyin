package db

import (
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"errors"
	"gorm.io/gorm"
)

type UserFavoriteVideo struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `gorm:"not null" json:"user_id"`
	VideoID   uint64 `gorm:"not null" json:"video_id"`
	IsDeleted uint8  `gorm:"default:0;not null" json:"is_deleted"`
}

func (n *UserFavoriteVideo) TableName() string {
	return constant.UserFavoriteVideosTableName
}

func FavoriteVideo(userID uint64, videoID uint64) error {
	if userID == 0 || videoID == 0 {
		return errno.UserRequestParameterError
	}

	userFavoriteVideo := &UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}

	// 先查询是否存在软删除的点赞数据
	result := global.DB.Where("is_deleted = ?", constant.DataDeleted).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 增加该视频的点赞数
		video := &Video{
			ID: videoID,
		}
		err := tx.First(&video).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&video).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
			return err
		}

		// 更新用户的点赞视频数
		u := &User{ID: userID}
		err = tx.Select("favorite_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("favorite_count", u.FavoriteCount+1).Error
		if err != nil {
			return err
		}
		if result.RowsAffected == 1 {
			// 如果有则修改为未删除
			return tx.Model(userFavoriteVideo).Update("is_deleted", constant.DataNotDeleted).Error
		}
		// 没有则新建
		return tx.Create(userFavoriteVideo).Error
	})
}

func CancelFavoriteVideo(userID uint64, videoID uint64) error {
	if userID == 0 || videoID == 0 {
		return errors.New("cancel favorite failed")
	}

	userFavoriteVideo := &UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}
	result := global.DB.Where("is_deleted = ?", constant.DataNotDeleted).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 增加该视频的点赞数
		video := &Video{
			ID: videoID,
		}
		err := tx.First(&video).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&video).Update("favorite_count", video.FavoriteCount-1).Error; err != nil {
			return err
		}

		// 更新用户的点赞视频数
		u := &User{ID: userID}
		err = tx.Select("favorite_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("favorite_count", u.FavoriteCount-1).Error
		if err != nil {
			return err
		}
		// 进行软删除
		return global.DB.Model(userFavoriteVideo).Update("is_deleted", constant.DataDeleted).Error
	})
}

func SelectFavoriteVideoListByUserID(toUserID uint64) ([]*Video, error) {
	res := make([]*Video, 0)
	// 使用子查询避免循环查询DB
	selectVideoIDSQL := global.DB.Select(`video_id`).
		Table("user_favorite_video").
		Where("user_id = ? AND is_deleted = ?", toUserID, constant.DataNotDeleted)
	global.DB.Where("id IN (?)", selectVideoIDSQL).Find(&res)
	return res, nil
}

func IsFavoriteVideo(userID, videoID uint64) bool {
	if userID == 0 || videoID == 0 {
		return false
	}
	ufv := make([]*UserFavoriteVideo, 0, 1)
	res := global.DB.Where("user_id = ? AND video_id = ? AND is_deleted = ?",
		userID, videoID, constant.DataNotDeleted).Limit(1).Find(&ufv)
	return res.RowsAffected == 1
}
