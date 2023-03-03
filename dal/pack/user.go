package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func User(u *db.User, isFollow bool) *api.User {
	if u == nil {
		hlog.Error("pack.user.User err:", errno.ServiceError)
		return nil
	}
	followCount := u.FollowingCount
	followerCount := u.FollowerCount
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
	if u == nil {
		hlog.Error("pack.user.UserInfo err:", errno.ServiceError)
		return nil
	}
	return &api.UserInfo{
		ID:              int64(u.ID),
		Name:            u.Username,
		FollowCount:     u.FollowingCount,
		FollowerCount:   u.FollowerCount,
		IsFollow:        isFollow,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  u.TotalFavorited,
		WorkCount:       u.WorkCount,
		FavoriteCount:   u.FavoriteCount,
	}
}

func FriendUser(u *db.User, isFollow bool, messageContent string, msgType uint8) *api.FriendUser {
	if u == nil {
		hlog.Error("pack.user.UserInfo err:", errno.ServiceError)
		return nil
	}
	followCount := u.FollowingCount
	followerCount := u.FollowerCount
	return &api.FriendUser{
		ID:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollow,
		Avatar:        u.Avatar,
		Message:       &messageContent,
		MsgType:       int8(msgType),
	}
}

func RelationData(data *db.RelationUserData) *api.User {
	if data == nil {
		return nil
	}
	followCount := int64(data.FollowingCount)
	followerCount := int64(data.FollowerCount)
	return &api.User{
		ID:            int64(data.UID),
		Name:          data.Username,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      data.IsFollow,
		Avatar:        data.Avatar,
	}
}

func RelationDataList(dataList []*db.RelationUserData) []*api.User {
	res := make([]*api.User, 0, len(dataList))
	for i := 0; i < len(dataList); i++ {
		res = append(res, RelationData(dataList[i]))
	}
	return res
}
