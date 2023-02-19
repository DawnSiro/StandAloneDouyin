package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func Comment(dbc *db.Comment, dbu *db.User, isFollow bool) *api.Comment {
	return &api.Comment{
		ID:         int64(dbc.ID),
		User:       User(dbu, isFollow),
		Content:    dbc.Content,
		CreateDate: dbc.CreatedTime.Format("01-02"), // 评论发布日期，格式 mm-dd
	}
}
