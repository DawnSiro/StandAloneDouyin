package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
)

func GetFeed(latestTime *int64, userID uint64) (*api.DouyinFeedResponse, error) {
	//// 方案一，每次都按单表查询，存在循环查询数据的问题，经过测试，开启协程进行异步也没有多少性能提升
	//videoList := make([]*api.Video, 0)
	//
	//videos, err := db.MGetVideos(constant.MaxVideoNum, latestTime)
	//if err != nil {
	//	hlog.Error("service.feed.GetFeed err:", err.Error())
	//	return nil, err
	//}
	//
	//// TODO 使用预处理等进行优化
	//for i := 0; i < len(videos); i++ {
	//	u, err := db.SelectUserByID(videos[i].AuthorID)
	//	if err != nil {
	//		hlog.Error("service.feed.GetFeed err:", err.Error())
	//		return nil, err
	//	}
	//
	//	var video *api.Video
	//	// 未登录默认未关注，未点赞
	//	if userID == 0 {
	//		video = pack.Video(videos[i], u,
	//			false, false)
	//	} else {
	//		video = pack.Video(videos[i], u,
	//			db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
	//	}
	//
	//	videoList = append(videoList, video)
	//}
	//
	//var nextTime *int64
	//// 没有视频的时候 nextTime 为 nil，会重置时间
	//if len(videos) != 0 {
	//	nextTime = new(int64)
	//	*nextTime = videos[len(videos)-1].PublishTime.UnixMilli()
	//}

	// 方案二，直接使用 JOIN 连接多个表数据，一次性查出所有数据。

	// Check if cache key is valid using Bloom filter
	cacheKey := fmt.Sprintf("feed:%d:%d", userID, latestTime)
	if bloomFilter.TestString(cacheKey) {
		cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
		if err == nil {
			// Cache hit, return cached feed data
			var cachedResponse api.DouyinFeedResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
				hlog.Error("service.feed.GetFeed err: Error decoding cached data, ", err.Error())
			} else {
				return &cachedResponse, nil
			}
		}
	}

	// Create a random duration for cache expiration
	minDuration := 5 * time.Hour
	maxDuration := 10 * time.Hour
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
		return GetFeed(latestTime, userID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	videoData, err := db.MSelectFeedVideoDataListByUserID(constant.MaxVideoNum, latestTime, userID)
	if err != nil {
		hlog.Error("service.feed.GetFeed err:", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	var nextTime *int64
	if len(videoData) != 0 {
		nextTime = new(int64)
		*nextTime, err = db.SelectPublishTimeByVideoID(videoData[len(videoData)-1].VID)
		if err != nil {
			hlog.Error("service.feed.GetFeed err:", err.Error())
			// Release cache status flag to allow other threads to update cache
			cacheMutex.Lock()
			delete(cacheStatus, cacheKey)
			cacheMutex.Unlock()
			return nil, err
		}
	}

	// Pack and return the response
	response := &api.DouyinFeedResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  pack.VideoDataList(videoData),
		NextTime:   nextTime,
	}

	// Store feed data in Redis cache with the random expiration time
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.feed.GetFeed err: Error storing data in cache, ", err.Error())
	}

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}
