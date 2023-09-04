package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"math/rand"
	"strconv"
	"sync"
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

	// Notify cache invalidation
	// Publish a message to the Redis channel indicating a friend list change
	global.UserInfoRC.Publish("friendList_changes", "friend_followed"+"&"+strconv.FormatUint(userID, 10)+"&"+strconv.FormatUint(toUserID, 10))

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

	// Notify cache invalidation
	// Publish a message to the Redis channel indicating a friend list change
	global.UserInfoRC.Publish("friendList_changes", "friend_unfollowed"+"&"+strconv.FormatUint(userID, 10)+"&"+strconv.FormatUint(toUserID, 10))

	return &api.DouyinRelationActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

// GetFollowList
// userID 为发送请求的用户ID，从 Token 里取到
// selectUserID 为需要查询的用户的ID，做为请求参数传递
func GetFollowList(userID, selectUserID uint64) (*api.DouyinRelationFollowListResponse, error) {
	// Check if cache key is valid using Bloom filter
	cacheKey := fmt.Sprintf("followList:%d", userID)
	if bloomFilter.TestString(cacheKey) {
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
	}

	// Create a random duration for cache expiration
	minDuration := 6 * time.Hour
	maxDuration := 12 * time.Hour
	cacheDuration := minDuration + time.Duration(rand.Intn(int(maxDuration-minDuration)))

	// Create a WaitGroup for the cache update operation
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Check if another thread is updating the cache
	cacheMutex.Lock()
	existingWaitGroup, exists := cacheStatus[cacheKey]
	if exists {
		cacheMutex.Unlock()
		existingWaitGroup.Wait()
		return GetFollowList(userID, selectUserID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	relationDataList, err := db.SelectFollowDataListByUserID(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowList err:", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinRelationFollowListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   pack.RelationDataList(relationDataList),
	}

	// Store the data in Redis cache with the random expiration time
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.relation.GetFollowList err: Error storing data in cache, ", err.Error())
	}

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}

func GetFollowerList(userID, selectUserID uint64) (*api.DouyinRelationFollowerListResponse, error) {
	cacheKey := fmt.Sprintf("followerList:%d", userID)
	if bloomFilter.TestString(cacheKey) {
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
	}

	// Create a random duration for cache expiration
	minDuration := 6 * time.Hour
	maxDuration := 12 * time.Hour
	cacheDuration := minDuration + time.Duration(rand.Intn(int(maxDuration-minDuration)))

	// Create a WaitGroup for the cache update operation
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Check if another thread is updating the cache
	cacheMutex.Lock()
	existingWaitGroup, exists := cacheStatus[cacheKey]
	if exists {
		cacheMutex.Unlock()
		existingWaitGroup.Wait()
		return GetFollowerList(userID, selectUserID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	relationDataList, err := db.SelectFollowerDataListByUserID(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err: ", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinRelationFollowerListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   pack.RelationDataList(relationDataList),
	}

	// Store the data in Redis cache with the random expiration time
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err: Error storing data in cache, ", err.Error())
	}

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}

func GetFriendList(userID uint64) (*api.DouyinRelationFriendListResponse, error) {
	cacheKey := fmt.Sprintf("friendList:%d", userID)
	if bloomFilter.TestString(cacheKey) {
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
	}

	// Create a random duration for cache expiration
	minDuration := 6 * time.Hour
	maxDuration := 12 * time.Hour
	cacheDuration := minDuration + time.Duration(rand.Intn(int(maxDuration-minDuration)))

	// Create a WaitGroup for the cache update operation
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Check if another thread is updating the cache
	cacheMutex.Lock()
	existingWaitGroup, exists := cacheStatus[cacheKey]
	if exists {
		cacheMutex.Unlock()
		existingWaitGroup.Wait()
		return GetFriendList(userID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	userList, err := db.GetFriendList(userID)
	if err != nil {
		hlog.Error("service.relation.GetFriendList err: ", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	// Convert to response format
	friendUserList := make([]*api.FriendUser, 0, len(userList))
	for _, u := range userList {
		msg, err := db.GetLatestMsg(userID, u.ID)
		if err != nil {
			hlog.Error("service.relation.GetFriendList err: ", err.Error())
			return nil, err
		}
		friendUserList = append(friendUserList, pack.FriendUser(u, db.IsFollow(userID, u.ID), msg.Content, msg.MsgType))
	}

	response := &api.DouyinRelationFriendListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   friendUserList,
	}

	// Store the data in Redis cache with the random expiration time
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.relation.GetFriendList err: Error storing data in cache, ", err.Error())
	}

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}
