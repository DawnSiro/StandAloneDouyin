package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetFeed(latestTime *int64, userID uint64) (*api.DouyinFeedResponse, error) {
	logTag := "service.feed.GetFeed err:"
	// 直接使用 JOIN 连接表数据，一次性查出所有数据。
	videoData, err := db.SelectFeedVideoDataListByUserID(constant.MaxVideoNum, latestTime)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	// 用户未登录，默认都是未关注未点赞，直接返回即可
	if userID == 0 {
		var nextTime *int64
		if len(videoData) != 0 {
			nextTime = new(int64)
			*nextTime, err = db.SelectPublishTimeByVideoID(videoData[len(videoData)-1].VID)
			if err != nil {
				hlog.Error(logTag, err.Error())
				return nil, err
			}
		}

		return &api.DouyinFeedResponse{
			StatusCode: errno.Success.ErrCode,
			VideoList:  pack.VideoDataList(videoData),
			NextTime:   nextTime,
		}, nil
	}

	// 获取用户的关注列表
	followSet, err := rdb.GetFollowUserIDSet(userID)
	if err != nil {
		idSet, err := db.SelectFollowUserIDSet(userID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
		for _, id := range idSet {
			followSet[id] = struct{}{}
		}
		err = rdb.SetFollowUserIDSet(userID, idSet)
		if err != nil {
			hlog.Error(logTag, err.Error())
		}
	}

	// 查询用户点赞列表
	favoriteVideoIDList, err := rdb.GetFavoriteVideoID(userID)
	favoriteVideoIDSet := make(map[uint64]struct{}, len(favoriteVideoIDList))
	for _, id := range favoriteVideoIDList {
		favoriteVideoIDSet[id] = struct{}{}
	}

	// 判断是否关注了视频作者及是否点赞了视频
	for _, video := range videoData {
		if _, ok := followSet[video.UID]; ok {
			video.IsFollow = true
		}
		if _, ok := favoriteVideoIDSet[video.VID]; ok {
			video.IsFavorite = true
		}
	}

	var nextTime *int64
	if len(videoData) != 0 {
		nextTime = new(int64)
		*nextTime, err = db.SelectPublishTimeByVideoID(videoData[len(videoData)-1].VID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
	}

	return &api.DouyinFeedResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  pack.VideoDataList(videoData),
		NextTime:   nextTime,
	}, nil
}
