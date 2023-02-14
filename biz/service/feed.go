package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"douyin/dal/pack"
)

func GetFeed(latestTime *int64, userID int64) (*api.DouyinFeedResponse, error) {
	res := new(api.DouyinFeedResponse)
	videoList := make([]*api.Video, 0)

	videos, err := db.MGetVideos(constant.MaxVideoNum, latestTime)
	if err != nil {
		return nil, err
	}

	// find author and pack data
	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(uint(videos[i].AuthorID))
		if err != nil {
			return nil, err
		}

		video, err := pack.Videos(videos[i], u,
			db.IsFollow(uint64(userID), uint64(u.ID)), db.IsFavorite(uint64(userID), uint64(videos[i].ID)))
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}

	res.VideoList = videoList

	return res, nil
}
