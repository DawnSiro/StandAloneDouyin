package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/constant"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetFeed(latestTime *int64, userID uint64) (*api.DouyinFeedResponse, error) {
	videoList := make([]*api.Video, 0)

	videos, err := db.MGetVideos(constant.MaxVideoNum, latestTime)
	if err != nil {
		hlog.Error("service.feed.GetFeed err:", err.Error())
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			hlog.Error("service.feed.GetFeed err:", err.Error())
			return nil, err
		}

		video := pack.Video(videos[i], u,
			db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
		videoList = append(videoList, video)
	}

	var nextTime *int64
	// 没有视频的时候 nextTime 为 nil，会重置时间
	if len(videos) != 0 {
		nextTime = new(int64)
		*nextTime = int64(videos[len(videos)-1].PublishTime.Second())
	}

	return &api.DouyinFeedResponse{
		StatusCode: 0,
		VideoList:  videoList,
		NextTime:   nextTime,
	}, nil
}
