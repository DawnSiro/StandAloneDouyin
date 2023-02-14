package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func Videos(v *db.Video, u *db.User, isFollow, isFavorite bool) (*api.Video, error) {
	followingCount := &u.FollowingCount
	followerCount := &u.FollowerCount
	author := &api.User{
		ID:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   followingCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
		Avatar:        u.Avatar,
	}
	res := &api.Video{
		ID:            int64(v.ID),
		Author:        author,
		PlayURL:       v.PlayURL,
		CoverURL:      v.CoverURL,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    isFavorite,
		Title:         v.Title,
	}

	return res, nil
}
