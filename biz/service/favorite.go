package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"github.com/go-redis/redis"
	"strconv"
)

// FavoriteAction this is a func for add Favorite or reduce Favorite
func FavoriteAction(req *api.DouyinFavoriteActionRequest) (api.DouyinFavoriteActionResponse, error) {
	var resp api.DouyinFavoriteActionResponse
	serverError := "服务器内部错误"
	serverOk := "ok"
	videoLikeKey := strconv.FormatInt(req.VideoID, 10) + "_video" + "_like"

	//TODO: token to userId
	//find like in redis is nil?
	likeCount, err := db.RDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		//find like_count in mysql
		likeInt64, err := db.SelectFavoriteCountByVideoId(req.VideoID)
		if err != nil {
			return api.DouyinFavoriteActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
			}, err
		}
		db.RDB.Set(videoLikeKey, likeInt64, 0)

	}
	if req.ActionType == 1 {
		//like
		//put it into redis
		likeInt64, err := strconv.ParseInt(likeCount, 10, 64)
		if err != nil {
			return api.DouyinFavoriteActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
			}, err
		}
		db.RDB.Set(videoLikeKey, likeInt64+1, 0)

		//TODO: token to userId
		//TODO: miss req.userId
		resultLike, err := db.Like(1, uint64(req.VideoID))
		if err != nil || resultLike == 0 {
			return api.DouyinFavoriteActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
			}, err
		}

	} else if req.ActionType == 2 {
		//cancel like
		//TODO: token to userId
		//put it into redis
		likeInt64, err := strconv.ParseInt(likeCount, 10, 64)
		if err != nil {
			return api.DouyinFavoriteActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError}, err
		}
		db.RDB.Set(videoLikeKey, likeInt64-1, 0)

		//TODO: token to userId
		//TODO: miss req.userId
		resultLike, err := db.UnLike(1, uint64(req.VideoID))
		if err != nil || resultLike == 0 {
			return api.DouyinFavoriteActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
			}, err
		}

	}

	resp.StatusCode = 0000
	resp.StatusMsg = &serverOk

	return resp, nil
}

// FavoriteList this is a func for get Favorite List
func FavoriteList(req *api.DouyinFavoriteListRequest) (api.DouyinFavoriteListResponse, error) {
	var resp api.DouyinFavoriteListResponse
	serverError := "服务器内部错误"
	serverOk := "ok"

	resultList, err := db.SelectFavoriteVideoListByUserId(uint64(req.UserID))
	if err != nil {
		return api.DouyinFavoriteListResponse{
			StatusCode: 1001,
			StatusMsg:  &serverError,
		}, err
	}

	resp.StatusCode = 0000
	resp.StatusMsg = &serverOk
	resp.VideoList = resultList

	return resp, nil
}
