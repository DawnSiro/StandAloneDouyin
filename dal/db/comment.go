package db

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"errors"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	VideoID uint64 `json:"video_id"`
	UserID  uint64 `json:"user_id"`
	Content string `json:"content"`
}

func (n *Comment) TableName() string {
	return constant.CommentTableName
}

func CreateComment(videoID uint64, content string, userID uint64) (*Comment, error) {
	comment := &Comment{
		Model:   gorm.Model{},
		VideoID: videoID,
		UserID:  userID,
		Content: content,
	}
	//
	if err := DB.Select("user_id", "video_id", "Content").Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func DeleteCommentByCommentID(commentID int64) (*Comment, error) {
	comment := &Comment{
		Model: gorm.Model{
			ID: uint(commentID),
		},
	}
	if err := DB.Delete(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func SelectCommentListByUserID(userID uint64, videoID uint64) ([]*api.Comment, error) {
	commentResult := new([]*Comment)
	err := DB.Where("video_id = ?", videoID).Find(&commentResult).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	results := make([]*api.Comment, 0)
	for i := 0; i < len(*commentResult); i++ {

		con1, err := SelectUserByUserID(uint((*commentResult)[i].UserID))
		if err != nil {
			return nil, err
		}
		con2, err := SelectAuthorIDByVideoID(int64(videoID))
		if err != nil || con2 == 0 {
			return nil, err
		}
		con1.IsFollow = IsFollow(userID, con2)
		commentTemp := &api.Comment{
			ID:         int64((*commentResult)[i].ID),
			User:       con1,
			Content:    (*commentResult)[i].Content,
			CreateDate: (*commentResult)[i].CreatedAt.String(),
		}
		results = append(results, commentTemp)
	}
	return results, nil
}

func IsCommentCreatedByMyself(userId uint64, commentId int64) (bool, error) {
	commentResult := &Comment{}
	result := DB.Where("user_id = ?", userId).Where("id = ?", commentId).Find(&commentResult)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}
	return true, nil
}
