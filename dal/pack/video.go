package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func Video(v *db.Video, u *db.User, isFollow, isFavorite bool) *api.Video {
	followingCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)
	author := &api.User{
		ID:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   &followingCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollow,
		Avatar:        u.Avatar,
	}

	res := &api.Video{
		ID:            int64(v.ID),
		Author:        author,
		PlayURL:       v.PlayURL,
		CoverURL:      v.CoverURL,
		FavoriteCount: int64(v.FavoriteCount),
		CommentCount:  int64(v.CommentCount),
		IsFavorite:    isFavorite,
		Title:         v.Title,
	}

	return res
}
