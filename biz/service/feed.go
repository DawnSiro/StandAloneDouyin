package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	// Check if feed data is available in Redis cache
	cacheKey := fmt.Sprintf("feed:%d:%d", userID, latestTime)
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

	// Cache miss, query the database
	videoData, err := db.MSelectFeedVideoDataListByUserID(constant.MaxVideoNum, latestTime, userID)
	if err != nil {
		hlog.Error("service.feed.GetFeed err:", err.Error())
		return nil, err
	}

	var nextTime *int64
	if len(videoData) != 0 {
		nextTime = new(int64)
		*nextTime, err = db.SelectPublishTimeByVideoID(videoData[len(videoData)-1].VID)
		if err != nil {
			hlog.Error("service.feed.GetFeed err:", err.Error())
			return nil, err
		}
	}

	// Pack and return the response
	response := &api.DouyinFeedResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  pack.VideoDataList(videoData),
		NextTime:   nextTime,
	}

	// Store feed data in Redis cache
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, 10*time.Minute).Err()
	if err != nil {
		hlog.Error("service.feed.GetFeed err: Error storing data in cache, ", err.Error())
	}

	return response, nil
}
