package db

import (
	"errors"

	"douyin/dal/model"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func CreateComment(comment *model.Comment) (*model.Comment, error) {
	// comment := &Comment{
	// 	VideoID:     videoID,
	// 	UserID:      userID,
	// 	Content:     content,
	// 	CreatedTime: time.Now(),
	// }
	// DB 层开事务来保证原子性
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 先查询 VideoID 是否存在，然后增加评论数，再创建评论
		video := &model.Video{
			ID: comment.VideoID,
		}
		err := tx.First(&video).Error
		if err != nil {
			return err
		}
		// 增加视频评论数
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Error
		if err != nil {
			return err
		}
		// 创建评论
		err = tx.Create(comment).Error
		if err != nil {
			return err
		}
		// 更新 Redis 缓存，如果返回错误会一起回滚，保证原子性和数据一致性
		return rdb.AddComment(video.ID, rdb.CommentInfo{
			ID:          comment.ID,
			VideoID:     comment.VideoID,
			UserID:      comment.UserID,
			Content:     comment.Content,
			CreatedTime: float64(comment.CreatedTime.UnixMilli()),
		})
	})
	if err != nil {
		hlog.Error("dal.db.comment.CreateComment err:", err.Error())
		return nil, err
	}

	return comment, nil
}

// DeleteCommentByID 通过评论ID 删除评论，默认使用软删除，提高性能
func DeleteCommentByID(videoID, commentID uint64) (*model.Comment, error) {
	comment := &model.Comment{
		ID: commentID,
	}
	// 先查询是否存在评论
	result := global.DB.Where("is_deleted = ?", constant.DataNotDeleted).Limit(1).Find(comment)
	if result.RowsAffected == 0 {
		return nil, errors.New("delete data failed")
	}

	// DB 层开事务来保证原子性
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 减少视频评论数
		video := &model.Video{
			ID: videoID,
		}
		err := tx.First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount-1).Error
		if err != nil {
			return err
		}
		// 删除评论
		err = tx.Model(comment).Update("is_deleted", constant.DataDeleted).Error
		if err != nil {
			return err
		}
		// 更新 Redis 缓存，如果返回错误会一起回滚，保证原子性和数据一致性
		return rdb.DeleteComment(videoID, commentID)
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func SelectCommentListByVideoID(videoID uint64) ([]*model.Comment, error) {
	res := make([]*model.Comment, 0)
	err := global.DB.Where("video_id = ? AND is_deleted = ?", videoID, constant.DataNotDeleted).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func IsCommentCreatedByMyself(userID uint64, commentID uint64) bool {
	result := global.DB.Where("id = ? AND user_id = ? AND is_deleted = ?", commentID, userID, constant.DataNotDeleted).
		Find(&model.Comment{})
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// SelectCommentDataByVideoIDAndUserID 查询评论数据，使用了JOIN，一次性查出所有的数据
func SelectCommentDataByVideoIDAndUserID(videoID, userID uint64) ([]*model.CommentData, error) {
	cs := make([]*model.CommentData, 0)
	err := global.DB.Select("c.id AS cid, c.content, c.created_time, "+
		"u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar, "+
		"IF(r.is_deleted = 0, TRUE, FALSE) AS is_follow").Table(constant.UserTableName+" AS u").
		Joins("LEFT JOIN "+constant.CommentTableName+" AS c ON u.id = c.user_id").
		Joins("LEFT JOIN "+constant.RelationTableName+" AS r ON u.id = r.`to_user_id` AND r.user_id = ?", userID).
		Where("`video_id` = ?", videoID).Scan(&cs).Error
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func SelectCommentDataByVideoID(videoID uint64) ([]model.Comment, error) {
	var comments []model.Comment

	if err := global.DB.Where("video_id = ?", videoID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
