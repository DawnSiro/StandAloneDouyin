package pack

import (
	"douyin/biz/model/api"
	"douyin/dal/model"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func User(u *model.User, isFollow bool) *api.User {
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

func UserInfo(u *model.User, isFollow bool) *api.UserInfo {
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

func CommentUser(u *model.User, isFollow bool) *api.CommentUser {
	if u == nil {
		hlog.Error("pack.user.CommentUser err:", errno.ServiceError)
		return nil
	}
	return &api.CommentUser{
		ID:       int64(u.ID),
		Name:     u.Username,
		IsFollow: isFollow,
		Avatar:   u.Avatar,
	}
}

func FriendUser(u *model.User, isFollow bool, messageContent string, msgType uint8) *api.FriendUser {
	if u == nil {
		hlog.Error("pack.user.UserInfo err:", errno.ServiceError)
		return nil
	}
	return &api.FriendUser{
		ID:       int64(u.ID),
		Name:     u.Username,
		IsFollow: isFollow,
		Avatar:   u.Avatar,
		Message:  &messageContent,
		MsgType:  int8(msgType),
	}
}

func FriendUserData(fuData *model.FriendUserData) *api.FriendUser {
	message := fuData.Message
	return &api.FriendUser{
		ID:       int64(fuData.ID),
		Name:     fuData.Name,
		IsFollow: fuData.IsFollow,
		Avatar:   fuData.Avatar,
		Message:  &message,
		MsgType:  fuData.MsgType,
	}
}

func FriendUserDataList(fuDataList []*model.FriendUserData) []*api.FriendUser {
	res := make([]*api.FriendUser, 0, len(fuDataList))
	for i := 0; i < len(fuDataList); i++ {
		res = append(res, FriendUserData(fuDataList[i]))
	}
	return res
}

func FollowUserWithRedis(data *model.FollowUserRedisData, isFollow bool) *api.FollowUser {
	if data == nil {
		hlog.Error("pack.user.FollowUserWithRedis err:", errno.ServiceError)
		return nil
	}
	return &api.FollowUser{
		ID:       int64(data.UID),
		Name:     data.Username,
		IsFollow: isFollow,
		Avatar:   data.Avatar,
	}
}

func FollowUser(data *model.FollowUserData, isFollow bool) *api.FollowUser {
	if data == nil {
		hlog.Error("pack.user.FollowUser err:", errno.ServiceError)
		return nil
	}
	return &api.FollowUser{
		ID:       int64(data.UID),
		Name:     data.Username,
		IsFollow: isFollow,
		Avatar:   data.Avatar,
	}
}

func FollowerUserWithRedis(data *model.FanUserRedisData, isFollow bool) *api.FollowerUser {
	if data == nil {
		hlog.Error("pack.user.FollowerUserWithRedis err:", errno.ServiceError)
		return nil
	}
	return &api.FollowerUser{
		ID:       int64(data.UID),
		Name:     data.Username,
		IsFollow: isFollow,
		Avatar:   data.Avatar,
	}
}

func FollowerUser(data *model.FanUserData, isFollow bool) *api.FollowerUser {
	if data == nil {
		hlog.Error("pack.user.FollowUser err:", errno.ServiceError)
		return nil
	}
	return &api.FollowerUser{
		ID:       int64(data.UID),
		Name:     data.Username,
		IsFollow: isFollow,
		Avatar:   data.Avatar,
	}
}

func RelationData(data *model.RelationUserData) *api.User {
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

func RelationDataList(dataList []*model.RelationUserData) []*api.User {
	res := make([]*api.User, 0, len(dataList))
	for i := 0; i < len(dataList); i++ {
		res = append(res, RelationData(dataList[i]))
	}
	return res
}
