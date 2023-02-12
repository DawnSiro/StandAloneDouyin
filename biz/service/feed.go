package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
)

func GetFeed(req *api.DouyinFeedRequest) (*api.DouyinFeedResponse, error) {
	res := new(api.DouyinFeedResponse)

	_, err := db.MGetVideos(constant.MaxVideoNum, req.LatestTime)
	if err != nil {
		return nil, err
	}

	// find author
	//db.SelectUserByUserID(videos)

	// find is favorite

	// pack

	return res, nil
}
