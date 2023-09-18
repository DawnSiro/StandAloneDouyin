package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/pulsar"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Follow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	logTag := "service.relation.Follow err:"
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己关注自己哦"
		hlog.Error(logTag, errNo.Error())
		return nil, errNo
	}

	// 查询缓存
	isFollow, err := rdb.IsFollow(userID, toUserID)
	// 缓存不存在则查询数据库
	if err != nil {
		hlog.Error(logTag, err.Error())
		isFollow = db.IsFollow(userID, toUserID)
		// 异步更新缓存，不阻塞
		go func() {
			set, err := db.SelectFollowUserIDSet(userID)
			if err != nil {
				hlog.Error(logTag, err.Error())
				return
			}
			err = rdb.SetFollowUserIDSet(userID, set)
			if err != nil {
				hlog.Error(logTag, err.Error())
			}
		}()
	}

	if isFollow {
		hlog.Error(logTag, errno.RepeatOperationError)
		return nil, errno.RepeatOperationError
	}

	// 关注操作
	// 放入消息队列
	err = pulsar.GetFollowActionMQInstance().FollowAction(toUserID, userID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	hlog.Debug("service.relation.Follow: publish a message")

	return &api.DouyinRelationActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func CancelFollow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	logTag := "service.relation.CancelFollow err:"
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己取关自己哦"
		hlog.Error(logTag, errNo.Error())
		return nil, errNo
	}

	// 查询缓存
	isFollow, err := rdb.IsFollow(userID, toUserID)
	// 缓存不存在则查询数据库
	if err != nil {
		hlog.Error(logTag, err.Error())
		isFollow = db.IsFollow(userID, toUserID)
		// 异步更新缓存，不阻塞
		go func() {
			set, err := db.SelectFollowUserIDSet(userID)
			if err != nil {
				hlog.Error(logTag, err.Error())
				return
			}
			err = rdb.SetFollowUserIDSet(userID, set)
			if err != nil {
				hlog.Error(logTag, err.Error())
			}
		}()
	}

	if !isFollow {
		hlog.Error(logTag, errno.RepeatOperationError)
		return nil, errno.RepeatOperationError
	}

	// 取消关注
	// 放入消息队列
	err = pulsar.GetFollowActionMQInstance().CancelFollowAction(toUserID, userID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	hlog.Debug("service.relation.CancelFollow: publish a message")

	return &api.DouyinRelationActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

// GetFollowList
// userID 为发送请求的用户ID，从 Token 里取到
// selectUserID 为需要查询的用户的ID，做为请求参数传递
func GetFollowList(userID, selectUserID uint64) (*api.DouyinRelationFollowListResponse, error) {
	logTag := "service.relation.GetFollowList err:"
	// 使用布隆过滤器判断用户ID是否存在
	if !global.UserIDBloomFilter.TestString(strconv.FormatUint(selectUserID, 10)) {
		hlog.Error(logTag, "布隆过滤器拦截")
		return nil, errno.UserRequestParameterError
	}

	followUserIDSet, err := rdb.GetFollowUserIDSet(userID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		// 查询数据库
		set, err := db.SelectFollowUserIDSet(userID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
		followUserIDSet = make(map[uint64]struct{}, len(set))
		for _, value := range set {
			followUserIDSet[value] = struct{}{}
		}
		// 设置缓存
		err = rdb.SetFollowUserIDSet(userID, set)
		if err != nil {
			hlog.Error(logTag, err.Error())
		}
	}

	// 获取缓存
	zSet, err := rdb.GetFollowUserZSet(selectUserID)
	// 查询缓存
	if err == nil {
		followUser := make([]*api.FollowUser, len(zSet))
		for i, data := range zSet {
			var isFollow bool
			// 判断是否在用户关注的集合中存在
			if _, ok := followUserIDSet[data.UID]; ok {
				isFollow = true
			}
			followUser[i] = pack.FollowUserWithRedis(data, isFollow)
		}

		return &api.DouyinRelationFollowListResponse{
			StatusCode: errno.Success.ErrCode,
			UserList:   followUser,
		}, nil
	}

	// 缓存未命中则查询数据库
	followUserList, err := db.SelectFollowUserListByUserID(selectUserID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	// 设置缓存
	err = rdb.SetFollowUserZSet(selectUserID, followUserList)
	if err != nil {
		hlog.Error(logTag, err.Error())
	}
	// 返回结果
	followUser := make([]*api.FollowUser, len(followUserList))
	for i, data := range followUserList {
		var isFollow bool
		// 判断是否在用户关注的集合中存在
		if _, ok := followUserIDSet[data.UID]; ok {
			isFollow = true
		}
		followUser[i] = pack.FollowUser(data, isFollow)
	}
	return &api.DouyinRelationFollowListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   followUser,
	}, nil
}

func GetFollowerList(userID, selectUserID uint64) (*api.DouyinRelationFollowerListResponse, error) {
	logTag := "service.relation.GetFollowerList err:"
	// 使用布隆过滤器判断用户ID是否存在
	if !global.UserIDBloomFilter.TestString(strconv.FormatUint(selectUserID, 10)) {
		hlog.Error(logTag, "布隆过滤器拦截")
		return nil, errno.UserRequestParameterError
	}

	// 查询用户关注列表缓存
	followUserIDSet, err := rdb.GetFollowUserIDSet(userID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		// 查询数据库
		set, err := db.SelectFollowUserIDSet(userID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
		followUserIDSet = make(map[uint64]struct{}, len(set))
		for _, value := range set {
			followUserIDSet[value] = struct{}{}
		}
		// 设置缓存
		err = rdb.SetFollowUserIDSet(userID, set)
		if err != nil {
			hlog.Error(logTag, err.Error())
		}
	}

	// 获取缓存
	zSet, err := rdb.GetFanUserZSet(selectUserID)
	// 查询缓存
	if err == nil {
		followUser := make([]*api.FollowerUser, len(zSet))
		for i, data := range zSet {
			var isFollow bool
			// 判断是否在用户关注的集合中存在
			if _, ok := followUserIDSet[data.UID]; ok {
				isFollow = true
			}
			followUser[i] = pack.FollowerUserWithRedis(data, isFollow)
		}

		return &api.DouyinRelationFollowerListResponse{
			StatusCode: errno.Success.ErrCode,
			UserList:   followUser,
		}, nil
	}

	// 缓存未命中则查询数据库
	followUserList, err := db.SelectFanUserListByUserID(selectUserID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	// 设置缓存
	err = rdb.SetFanUserZSet(selectUserID, followUserList)
	if err != nil {
		hlog.Error(logTag, err.Error())
	}
	// 返回结果
	followUser := make([]*api.FollowerUser, len(followUserList))
	for i, data := range followUserList {
		var isFollow bool
		// 判断是否在用户关注的集合中存在
		if _, ok := followUserIDSet[data.UID]; ok {
			isFollow = true
		}
		followUser[i] = pack.FollowerUser(data, isFollow)
	}
	return &api.DouyinRelationFollowerListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   followUser,
	}, nil
}

func GetFriendList(userID uint64) (*api.DouyinRelationFriendListResponse, error) {
	logTag := "service.relation.GetFriendList err: "
	// db 层进行了处理，解决了循环查询 DB 的问题
	fuDataList, err := db.SelectFriendDataList(userID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	return &api.DouyinRelationFriendListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   pack.FriendUserDataList(fuDataList),
	}, nil
}
