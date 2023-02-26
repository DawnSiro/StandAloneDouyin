package db

import (
	"douyin/pkg/global"
	"errors"
	"time"

	"douyin/pkg/constant"
	"douyin/pkg/errno"

	"gorm.io/gorm"
)

type Comment struct {
	ID          uint64    `json:"id"`
	IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	VideoID     uint64    `gorm:"not null" json:"video_id"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	Content     string    `gorm:"type:varchar(255);not null" json:"content"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
}

func (n *Comment) TableName() string {
	return constant.CommentTableName
}

func CreateComment(videoID uint64, content string, userID uint64) (*Comment, error) {
	comment := &Comment{
		VideoID:     videoID,
		UserID:      userID,
		Content:     content,
		CreatedTime: time.Now(),
	}
	// DB 层开事务来保证原子性
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建评论
		err := global.DB.Create(comment).Error
		if err != nil {
			return err
		}
		// 增加视频评论数
		video := &Video{
			ID: videoID,
		}
		err = global.DB.First(&video).Error
		if err != nil {
			return err
		}
		return global.DB.Model(&video).Update("comment_count", video.CommentCount+1).Error
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteCommentByID 通过评论ID 删除评论，默认使用软删除，提高性能
func DeleteCommentByID(videoID, commentID uint64) (*Comment, error) {
	comment := &Comment{
		ID: commentID,
	}
	// 先查询是否存在评论
	result := global.DB.Where("is_deleted = ?", constant.DataNotDeleted).Limit(1).Find(comment)
	if result.RowsAffected == 0 {
		return nil, errors.New("delete data failed")
	}

	// DB 层开事务来保证原子性
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除评论
		result = global.DB.Model(comment).Update("is_deleted", constant.DataDeleted)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errno.UserRequestParameterError
		}
		// 减少视频评论数
		video := &Video{
			ID: videoID,
		}
		err := global.DB.First(&video).Error
		if err != nil {
			return err
		}
		return global.DB.Model(&video).Update("comment_count", video.CommentCount-1).Error
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func SelectCommentListByVideoID(videoID uint64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	err := global.DB.Where("video_id = ? AND is_deleted = ?", videoID, constant.DataNotDeleted).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func IsCommentCreatedByMyself(userID uint64, commentID uint64) bool {
	result := global.DB.Where("id = ? AND user_id = ? AND is_deleted = ?", commentID, userID, constant.DataNotDeleted).
		Find(&Comment{})
	if result.RowsAffected == 0 {
		return false
	}
	return true
}
