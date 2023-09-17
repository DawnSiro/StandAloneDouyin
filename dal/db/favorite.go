package db

import (
	"errors"
	"time"

	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"gorm.io/gorm"
)

func FavoriteVideo(userID uint64, videoID uint64) error {
	if userID == 0 || videoID == 0 {
		return errno.UserRequestParameterError
	}

	userFavoriteVideo := &model.UserFavoriteVideo{
		UserID:      userID,
		VideoID:     videoID,
		CreatedTime: time.Now(),
	}

	// 先查询是否存在软删除的点赞数据
	result := global.DB.Where("user_id = ? AND video_id = ? AND is_deleted = ?",
		userID, videoID, constant.DataDeleted).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 增加该视频的点赞数
		video := &model.Video{
			ID: videoID,
		}
		err := tx.Select("favorite_count, author_id").First(&video).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&video).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
			return err
		}

		// 更新用户的点赞视频数
		u := &model.User{ID: userID}
		err = tx.Select("favorite_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("favorite_count", u.FavoriteCount+1).Error
		if err != nil {
			return err
		}

		// 更新被点赞用户的总获赞数
		author := &model.User{ID: video.AuthorID}
		err = tx.Select("total_favorited").First(author).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("total_favorited", author.TotalFavorited+1).Error
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

	userFavoriteVideo := &model.UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}
	result := global.DB.Where("user_id = ? AND video_id = ?",
		userID, videoID).Limit(1).Find(userFavoriteVideo)
	if result.Error != nil {
		return result.Error
	}
	if userFavoriteVideo.IsDeleted == constant.DataDeleted {
		return errors.New("用户重复操作")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 增加该视频的点赞数
		video := &model.Video{
			ID: videoID,
		}
		err := tx.Select("favorite_count, author_id").First(&video).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&video).Update("favorite_count", video.FavoriteCount-1).Error; err != nil {
			return err
		}

		// 更新用户的点赞视频数
		u := &model.User{ID: userID}
		err = tx.Select("favorite_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("favorite_count", u.FavoriteCount-1).Error
		if err != nil {
			return err
		}

		// 更新被点赞用户的总获赞数
		author := &model.User{ID: video.AuthorID}
		err = tx.Select("total_favorited").First(author).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("total_favorited", author.TotalFavorited-1).Error
		if err != nil {
			return err
		}

		// 进行软删除
		return global.DB.Model(userFavoriteVideo).Update("is_deleted", constant.DataDeleted).Error
	})
}

func SelectFavoriteVideoListByUserID(toUserID uint64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	// 使用子查询避免循环查询DB
	selectVideoIDSQL := global.DB.Select(`video_id`).
		Table(constant.UserFavoriteVideosTableName).
		Where("user_id = ? AND is_deleted = ?", toUserID, constant.DataNotDeleted)
	global.DB.Where("id IN (?)", selectVideoIDSQL).Find(&res)
	return res, nil
}

func IsFavoriteVideo(userID, videoID uint64) bool {
	if userID == 0 || videoID == 0 {
		return false
	}
	ufv := make([]*model.UserFavoriteVideo, 0, 1)
	res := global.DB.Where("user_id = ? AND video_id = ? AND is_deleted = ?",
		userID, videoID, constant.DataNotDeleted).Limit(1).Find(&ufv)
	return res.RowsAffected == 1
}

func SelectTopFavoriteVideos(limit int) ([]model.Video, error) {
	var videos []model.Video

	if err := global.DB.Order("favorite_count desc").Limit(limit).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}
