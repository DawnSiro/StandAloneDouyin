package db

import (
	"douyin/biz/model/api"
	"douyin/constants"
	"errors"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	VideoId uint64 `json:"video_id"`
	UserId  uint64 `json:"user_id"`
	Content string `json:"content"`
}

func (n *Comment) TableName() string {
	return constants.CommentTableName
}

func CreateCommentByUserIdAndVideoIdAndContent(req api.DouyinCommentActionRequest) (*Comment, error) {
	comment := &Comment{
		Model:   gorm.Model{},
		VideoId: uint64(req.VideoID),
		UserId:  uint64(req.UserID),
		Content: *req.CommentText,
	}
	//
	if err := DB.Select("UserId", "VideoId", "Content").Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func DeleteCommentByCommentId(commentId int64) (*Comment, error) {
	comment := &Comment{
		Model: gorm.Model{
			ID: uint(commentId),
		},
	}
	if err := DB.Delete(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func SelectCommentListByUserId(userId uint64, videoId uint64) ([]*api.Comment, error) {
	commentResult := new([]*Comment)
	result := DB.Where("user_id", userId).Where("video_id", videoId).Find(&commentResult)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	results := make([]*api.Comment, 0)
	for i := 0; i < len(*commentResult); i++ {

		con1, err := SelectUserByUserId(uint(userId))
		if err != nil {
			return nil, nil
		}
		con2, err := SelectAuthorIdByVideoId(int64(videoId))
		if err != nil || con2 == 0 {
			return nil, nil
		}
		con1.IsFollow = IsFollow(userId, con2)
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
