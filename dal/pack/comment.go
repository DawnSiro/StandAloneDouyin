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
