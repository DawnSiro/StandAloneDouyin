package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/constant"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetFeed(latestTime *int64, userID uint64) (*api.DouyinFeedResponse, error) {
	videoList := make([]*api.Video, 0)

	videos, err := db.MGetVideos(constant.MaxVideoNum, latestTime)
	if err != nil {
		hlog.Error("service.feed.GetFeed err:", err.Error())
		return nil, err
	}

	// TODO 使用预处理等进行优化
	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			hlog.Error("service.feed.GetFeed err:", err.Error())
			return nil, err
		}

		var video *api.Video
		// 未登录默认未关注，未点赞
		if userID == 0 {
			video = pack.Video(videos[i], u,
				false, false)
		} else {
			video = pack.Video(videos[i], u,
				db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
		}

		videoList = append(videoList, video)
	}

	var nextTime *int64
	// 没有视频的时候 nextTime 为 nil，会重置时间
	if len(videos) != 0 {
		nextTime = new(int64)
		*nextTime = videos[len(videos)-1].PublishTime.UnixMilli()
	}

	return &api.DouyinFeedResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  videoList,
		NextTime:   nextTime,
	}, nil
}
