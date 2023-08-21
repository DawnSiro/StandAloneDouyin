package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Follow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己关注自己哦"
		hlog.Error("service.relation.Follow err:", errNo.Error())
		return nil, errNo
	}
	isFollow := db.IsFollow(userID, toUserID)
	if isFollow {
		hlog.Error("service.relation.Follow err:", errno.RepeatOperationError)
		return nil, errno.RepeatOperationError
	}

	//关注操作
	err := db.Follow(userID, toUserID)
	if err != nil {
		hlog.Error("service.relation.Follow err:", err.Error())
		return nil, err
	}
	return &api.DouyinRelationActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func CancelFollow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己取关自己哦"
		hlog.Error("service.relation.CancelFollow err:", errNo.Error())
		return nil, errNo
	}
	//取消关注
	err := db.CancelFollow(userID, toUserID)
	if err != nil {
		hlog.Error("service.relation.CancelFollow err:", err.Error())
		return nil, err
	}
	return &api.DouyinRelationActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

// GetFollowList
// userID 为发送请求的用户ID，从 Token 里取到
// selectUserID 为需要查询的用户的ID，做为请求参数传递
func GetFollowList(userID, selectUserID uint64) (*api.DouyinRelationFollowListResponse, error) {
	//dbUserList, err := db.GetFollowList(selectUserID)
	//if err != nil {
	//	hlog.Error("service.relation.GetFollowList err:", err.Error())
	//	return nil, err
	//}
	//
	//// 提前申请好数组大小来避免后续扩容
	//userList := make([]*api.User, 0, len(dbUserList))
	//// TODO 存在循环查询DB
	//for _, v := range dbUserList {
	//	if userID == selectUserID {
	//		// 自己的关注列表自己当然都关注了，无需查数据库
	//		userList = append(userList, pack.User(v, true))
	//	} else {
	//		// 这里要查的是，自己是否关注了查询的用户的关注列表的人
	//		userList = append(userList, pack.User(v, db.IsFollow(userID, v.ID)))
	//	}
	//}

	// Check if data is available in Redis
	cacheKey := fmt.Sprintf("followlist:%d", userID)
	cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
	if err == nil {
		// Cache hit, return cached data
		var cachedResponse api.DouyinRelationFollowListResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
			hlog.Error("service.relation.GetFollowList err: Error decoding cached data, ", err.Error())
		} else {
			return &cachedResponse, nil
		}
	}

	// Cache miss, query the database
	relationDataList, err := db.SelectFollowDataListByUserID(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowList err:", err.Error())
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinRelationFollowListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   pack.RelationDataList(relationDataList),
	}

	// Store the data in Redis cache
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, 1*time.Hour).Err()
	if err != nil {
		hlog.Error("service.relation.GetFollowList err: Error storing data in cache, ", err.Error())
	}

	return response, nil

}

func GetFollowerList(userID, selectUserID uint64) (*api.DouyinRelationFollowerListResponse, error) {
	//dbUserList, err := db.GetFollowerList(selectUserID)
	//if err != nil {
	//	hlog.Error("service.relation.GetFollowerList err:", err.Error())
	//	return nil, err
	//}
	//
	//// 提前申请好数组大小来避免后续扩容
	//userList := make([]*api.User, 0, len(dbUserList))
	//// TODO 存在循环查询DB
	//for _, v := range dbUserList {
	//	// 这里要查的是，自己是否关注了查询的用户的粉丝列表的人
	//	userList = append(userList, pack.User(v, db.IsFollow(userID, v.ID)))
	//}

	// Check if data is available in Redis
	cacheKey := fmt.Sprintf("followerlist:%d", userID)
	cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
	if err == nil {
		// Cache hit, return cached data
		var cachedResponse api.DouyinRelationFollowerListResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
			hlog.Error("service.relation.GetFollowerList err: Error decoding cached data, ", err.Error())
		} else {
			return &cachedResponse, nil
		}
	}

	// Cache miss, query the database
	relationDataList, err := db.SelectFollowerDataListByUserID(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err: ", err.Error())
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinRelationFollowerListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   pack.RelationDataList(relationDataList),
	}

	// Store the data in Redis cache
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, 1*time.Hour).Err()
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err: Error storing data in cache, ", err.Error())
	}

	return response, nil

}

func GetFriendList(userID uint64) (*api.DouyinRelationFriendListResponse, error) {
	// Check if data is available in Redis
	cacheKey := fmt.Sprintf("friendlist:%d", userID)
	cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
	if err == nil {
		// Cache hit, return cached data
		var cachedResponse api.DouyinRelationFriendListResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
			hlog.Error("service.relation.GetFriendList err: Error decoding cached data, ", err.Error())
		} else {
			return &cachedResponse, nil
		}
	}

	// Cache miss, query the database
	userList, err := db.GetFriendList(userID)
	if err != nil {
		hlog.Error("service.relation.GetFriendList err: ", err.Error())
		return nil, err
	}

	// TODO 存在循环查询DB
	friendUserList := make([]*api.FriendUser, 0, len(userList))
	for _, u := range userList {
		msg, err := db.GetLatestMsg(userID, u.ID)
		if err != nil {
			hlog.Error("service.relation.GetFriendList err: ", err.Error())
			return nil, err
		}
		friendUserList = append(friendUserList, pack.FriendUser(u, db.IsFollow(userID, u.ID), msg.Content, msg.MsgType))
	}

	// Convert to response format
	response := &api.DouyinRelationFriendListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   friendUserList,
	}

	// Store the data in Redis cache
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, 1*time.Hour).Err()
	if err != nil {
		hlog.Error("service.relation.GetFriendList err: Error storing data in cache, ", err.Error())
	}

	return response, nil
}
