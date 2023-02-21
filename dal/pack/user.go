package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

func User(u *db.User, isFollow bool) *api.User {
	followCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)
	return &api.User{
		ID:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollow,
		Avatar:        u.Avatar,
	}
}

func UserInfo(u *db.User, isFollow bool) *api.UserInfo {
	return &api.UserInfo{
		ID:              int64(u.ID),
		Name:            u.Username,
		FollowCount:     int64(u.FollowingCount),
		FollowerCount:   int64(u.FollowerCount),
		IsFollow:        isFollow,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		WorkCount:       int64(u.WorkCount),
		FavoriteCount:   int64(u.FavoriteCount),
	}
}
