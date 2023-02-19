package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"github.com/go-redis/redis"
	"strconv"
)

// FavoriteVideo this is a func for add Favorite or reduce Favorite
func FavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	videoLikeKey := strconv.FormatUint(videoID, 10) + "_video" + "_like"

	likeCount, err := db.RDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			return nil, err
		}
		db.RDB.Set(videoLikeKey, likeInt64, 0)
	}
	//like
	//put it into redis
	likeInt64, err := strconv.ParseUint(likeCount, 10, 64)
	if err != nil {
		return nil, err
	}

	err = db.FavoriteVideo(userID, videoID)
	if err != nil {
		return nil, err
	}
	db.RDB.Set(videoLikeKey, likeInt64+1, 0)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  nil,
	}, nil
}

func CancelFavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	videoLikeKey := strconv.FormatUint(videoID, 10) + "_video" + "_like"

	likeCount, err := db.RDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		//find like_count in mysql
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			return nil, err
		}
		db.RDB.Set(videoLikeKey, likeInt64, 0)
	}
	//cancel like
	//put it into redis
	likeInt64, err := strconv.ParseInt(likeCount, 10, 64)
	if err != nil {
		return nil, err
	}

	err = db.CancelFavoriteVideo(userID, videoID)
	if err != nil {
		return nil, err
	}
	db.RDB.Set(videoLikeKey, likeInt64-1, 0)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  nil,
	}, nil
}

func FavoriteList(userID, selectUserID uint64) (*api.DouyinFavoriteListResponse, error) {
	videos, err := db.SelectFavoriteVideoListByUserID(selectUserID)
	if err != nil {
		return nil, err
	}

	videoList := make([]*api.Video, 0)
	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			return nil, err
		}
		video := pack.Video(videos[i], u,
			db.IsFollow(userID, selectUserID), db.IsFavoriteVideo(userID, videos[i].ID))
		videoList = append(videoList, video)
	}

	return &api.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		VideoList:  videoList,
	}, nil
}
