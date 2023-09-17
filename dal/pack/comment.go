package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/model"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Comment(dbc *model.Comment, dbu *model.User, isFollow bool) *api.Comment {
	if dbc == nil || dbu == nil {
		hlog.Error("pack.comment.Comment err:", errno.ServiceError)
		return nil
	}

	return &api.Comment{
		ID:         int64(dbc.ID),
		User:       CommentUser(dbu, isFollow),
		Content:    dbc.Content,
		CreateDate: dbc.CreatedTime.Format("01-02"), // 评论发布日期，格式 mm-dd
	}
}

func CommentData(data *model.CommentData) *api.Comment {
	if data == nil {
		hlog.Error("pack.comment.CommentData err:", errno.ServiceError)
		return nil
	}
	u := &api.CommentUser{
		ID:       int64(data.UID),
		Name:     data.Username,
		IsFollow: data.IsFollow,
		Avatar:   data.Avatar,
	}
	return &api.Comment{
		ID:         int64(data.CID),
		User:       u,
		Content:    data.Content,
		CreateDate: data.CreatedTime.Format("01-02"), // 评论发布日期，格式 mm-dd
	}
}

func CommentDataList(cdList []*model.CommentData) []*api.Comment {
	res := make([]*api.Comment, 0, len(cdList))
	for i := 0; i < len(cdList); i++ {
		res = append(res, CommentData(cdList[i]))
	}
	return res
}
