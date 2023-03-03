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

func VideoData(data *db.VideoData) *api.Video {
	if data == nil {
		return nil
	}
	followCount := data.FollowCount
	followerCount := data.FollowerCount
	author := &api.UserInfo{
		ID:              int64(data.UID),
		Name:            data.Username,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        data.IsFollow,
		Avatar:          data.Avatar,
		BackgroundImage: data.BackgroundImage,
		Signature:       data.Signature,
		TotalFavorited:  data.TotalFavorited,
		WorkCount:       data.WorkCount,
		FavoriteCount:   data.UserFavoriteCount,
	}
	return &api.Video{
		ID:            int64(data.VID),
		Author:        author,
		PlayURL:       data.PlayURL,
		CoverURL:      data.CoverURL,
		FavoriteCount: data.FavoriteCount,
		CommentCount:  data.CommentCount,
		IsFavorite:    data.IsFavorite,
		Title:         data.Title,
	}
}

func VideoDataList(dataList []*db.VideoData) []*api.Video {
	res := make([]*api.Video, 0, len(dataList))
	for i := 0; i < len(dataList); i++ {
		res = append(res, VideoData(dataList[i]))
	}
	return res
}
