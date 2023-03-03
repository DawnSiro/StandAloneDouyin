package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Comment(dbc *db.Comment, dbu *db.User, isFollow bool) *api.Comment {
	if dbc == nil || dbu == nil {
		hlog.Error("pack.comment.Comment err:", errno.ServiceError)
		return nil
	}

	return &api.Comment{
		ID:         int64(dbc.ID),
		User:       User(dbu, isFollow),
		Content:    dbc.Content,
		CreateDate: dbc.CreatedTime.Format("01-02"), // 评论发布日期，格式 mm-dd
	}
}

func CommentData(data *db.CommentData) *api.Comment {
	if data == nil {
		return nil
	}
	followCount := int64(data.FollowingCount)
	followerCount := int64(data.FollowerCount)
	u := &api.User{
		ID:            int64(data.UID),
		Name:          data.Username,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      data.IsFollow,
		Avatar:        data.Avatar,
	}
	return &api.Comment{
		ID:         int64(data.CID),
		User:       u,
		Content:    data.Content,
		CreateDate: data.CreatedTime.Format("01-02"), // 评论发布日期，格式 mm-dd
	}
}

func CommentDataList(cdList []*db.CommentData) []*api.Comment {
	res := make([]*api.Comment, 0, len(cdList))
	for i := 0; i < len(cdList); i++ {
		res = append(res, CommentData(cdList[i]))
	}
	return res
}
