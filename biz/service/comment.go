package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/util/sensitive"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	// 删除redis评论列表缓存
	// 使用 strings.Builder 来优化字符串的拼接
	//var builder strings.Builder
	//builder.WriteString(strconv.FormatUint(videoID, 10))
	//builder.WriteString("_video_comments")
	//delCommentListKey := builder.String()
	//hlog.Info("service.comment.PostComment delCommentListKey:", delCommentListKey)

	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errno.ContainsProhibitedSensitiveWordsError
	}

	dbc, err := db.CreateComment(videoID, commentText, userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}

	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	// 查询此评论是否是本人发送的
	isComment := db.IsCommentCreatedByMyself(userID, commentID)
	// 非本人评论
	if !isComment {
		hlog.Error("service.comment.DeleteComment err:", errno.DeletePermissionError)
		return nil, errno.DeletePermissionError
	}

	dbc, err := db.DeleteCommentByID(videoID, commentID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func GetCommentList(userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	commentData, err := db.SelectCommentDataByVideoIDANDUserID(videoID, userID)
	if err != nil {
		return nil, err
	}
	return &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: pack.CommentDataList(commentData),
	}, nil
}
