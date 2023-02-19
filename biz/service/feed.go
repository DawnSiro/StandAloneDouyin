package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/constant"
)

func GetFeed(latestTime *int64, userID uint64) (*api.DouyinFeedResponse, error) {
	videoList := make([]*api.Video, 0)

	videos, err := db.MGetVideos(constant.MaxVideoNum, latestTime)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			return nil, err
		}

		video := pack.Video(videos[i], u,
			db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
		videoList = append(videoList, video)
	}

	nextTime := int64(videos[len(videos)-1].PublishTime.Second())

	return &api.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		VideoList:  videoList,
		NextTime:   &nextTime,
	}, nil
}
