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
	followCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)
	var workCount = int64(u.WorkCount)
	var favoriteCount = int64(u.FavoriteCount)
	return &api.UserInfo{
		ID:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollow,
		Avatar:        u.Avatar,
		WorkCount:     &workCount,
		FavoriteCount: &favoriteCount,
	}
}
