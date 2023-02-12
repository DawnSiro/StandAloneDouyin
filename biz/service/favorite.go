package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis"
	"strconv"
)

// FavoriteAction this is a func for add Favorite or reduce Favorite
func FavoriteAction(req *api.DouyinFavoriteActionRequest, c *app.RequestContext) (api.DouyinFavoriteActionResponse, error) {
	var resp api.DouyinFavoriteActionResponse
	videoLikeKey := strconv.FormatInt(req.VideoID, 10) + "_video" + "_like"

	//find like in redis is nil?
	userId := c.GetInt64(constant.IdentityKey)
	likeCount, err := db.RDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		//find like_count in mysql
		likeInt64, err := db.SelectFavoriteCountByVideoId(req.VideoID)
		if err != nil {
			return api.DouyinFavoriteActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
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
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, err
		}
		db.RDB.Set(videoLikeKey, likeInt64+1, 0)

		resultLike, err := db.Like(uint64(userId), uint64(req.VideoID))
		if err != nil || resultLike == 0 {
			return api.DouyinFavoriteActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, err
		}

	} else if req.ActionType == 2 {
		//cancel like
		//put it into redis
		likeInt64, err := strconv.ParseInt(likeCount, 10, 64)
		if err != nil {
			return api.DouyinFavoriteActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, err
		}
		db.RDB.Set(videoLikeKey, likeInt64-1, 0)

		resultLike, err := db.UnLike(uint64(userId), uint64(req.VideoID))
		if err != nil || resultLike == 0 {
			return api.DouyinFavoriteActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, err
		}

	}

	resp.StatusCode = int64(api.ErrCode_SuccessCode)
	return resp, nil
}

// FavoriteList this is a func for get Favorite List
func FavoriteList(req *api.DouyinFavoriteListRequest, c *app.RequestContext) (api.DouyinFavoriteListResponse, error) {
	var resp api.DouyinFavoriteListResponse

	userId := c.GetInt64(constant.IdentityKey)
	resultList, err := db.SelectFavoriteVideoListByUserId(uint64(userId), uint64(req.UserID))
	if err != nil {
		return api.DouyinFavoriteListResponse{
			StatusCode: int64(api.ErrCode_ServiceErrCode),
		}, err
	}

	resp.StatusCode = int64(api.ErrCode_SuccessCode)
	resp.VideoList = resultList

	return resp, nil
}