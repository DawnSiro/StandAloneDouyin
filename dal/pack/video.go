package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Video(v *db.Video, u *db.User, isFollow, isFavorite bool) *api.Video {
	if v == nil || u == nil {
		hlog.Error("pack.video.Video err:", errno.ServiceError)
		return nil
	}
	author := &api.UserInfo{
		ID:              int64(u.ID),
		Name:            u.Username,
		FollowCount:     int64(u.FollowingCount),
		FollowerCount:   int64(u.FollowerCount),
		IsFollow:        isFollow,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  int64(u.TotalFavorited),
		WorkCount:       int64(u.WorkCount),
		FavoriteCount:   int64(u.FavoriteCount),
	}
	return &api.Video{
		ID:            int64(v.ID),
		Author:        author,
		PlayURL:       v.PlayURL,
		CoverURL:      v.CoverURL,
		FavoriteCount: int64(v.FavoriteCount),
		CommentCount:  int64(v.CommentCount),
		IsFavorite:    isFavorite,
		Title:         v.Title,
	}
}
