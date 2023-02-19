package db

import (
	"douyin/pkg/constant"
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
		return errors.New("favorite failed")
	}

	userFavoriteVideo := &UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}

	// 先查询是否存在软删除的点赞数据
	result := DB.Where("is_deleted = ?", constant.DataDeleted).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		// 先更新用户的点赞视频数
		u := &User{ID: userID}
		err := tx.Select("favorite_count").First(u).Error
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
	result := DB.Where("is_deleted = ?", constant.DataNotDeleted).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户的点赞视频数
		u := &User{ID: userID}
		err := tx.Select("favorite_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("favorite_count", u.FavoriteCount-1).Error
		if err != nil {
			return err
		}
		// 进行软删除
		return DB.Model(userFavoriteVideo).Update("is_deleted", constant.DataDeleted).Error
	})
}

func SelectFavoriteVideoListByUserID(toUserID uint64) ([]*Video, error) {
	res := make([]*Video, 0)

	err := DB.Where("author_id = ?", toUserID).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func IsFavoriteVideo(userID, videoID uint64) bool {
	if userID == 0 || videoID == 0 {
		return false
	}
	ufv := make([]*UserFavoriteVideo, 1)
	res := DB.Where("user_id = ? AND video_id = ? AND is_deleted = ?",
		userID, videoID, constant.DataNotDeleted).Limit(1).Find(&ufv)
	return res.RowsAffected == 1
}
