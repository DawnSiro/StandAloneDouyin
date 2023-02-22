package service

import (
	"strconv"
	"strings"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

// FavoriteVideo this is a func for add Favorite or reduce Favorite
func FavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_like")
	videoLikeKey := builder.String()

	likeCount, err := db.VideoFRDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
			return nil, err
		}
		db.VideoFRDB.Set(videoLikeKey, likeInt64, 0)
	}
	//like
	//put it into redis
	likeInt64, err := strconv.ParseUint(likeCount, 10, 64)
	if err != nil {
		hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
		return nil, err
	}

	err = db.FavoriteVideo(userID, videoID)
	if err != nil {
		hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
		return nil, err
	}
	// 如果 DB 层事务回滚了，err 就不为 nil，Redis 里的数据就不会更新
	db.VideoFRDB.Set(videoLikeKey, likeInt64+1, 0)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: 0,
	}, nil
}

func CancelFavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_like")
	videoLikeKey := builder.String()

	likeCount, err := db.VideoFRDB.Get(videoLikeKey).Result()
	if err == redis.Nil {
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
			return nil, err
		}
		db.VideoFRDB.Set(videoLikeKey, likeInt64, 0)
	}

	likeInt64, err := strconv.ParseInt(likeCount, 10, 64)
	if err != nil {
		hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
		return nil, err
	}

	err = db.CancelFavoriteVideo(userID, videoID)
	if err != nil {
		hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
		return nil, err
	}
	// 如果 DB 层事务回滚了，err 就不为 nil，Redis 里的数据就不会更新
	db.VideoFRDB.Set(videoLikeKey, likeInt64-1, 0)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: 0,
	}, nil
}

func FavoriteList(userID, selectUserID uint64) (*api.DouyinFavoriteListResponse, error) {
	videos, err := db.SelectFavoriteVideoListByUserID(selectUserID)
	if err != nil {
		hlog.Error("service.favorite.FavoriteList err:", err.Error())
		return nil, err
	}

	// TODO 优化循环查询数据库问题
	videoList := make([]*api.Video, 0)
	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			hlog.Error("service.favorite.FavoriteList err:", err.Error())
			return nil, err
		}
		video := pack.Video(videos[i], u,
			db.IsFollow(userID, selectUserID), db.IsFavoriteVideo(userID, videos[i].ID))
		videoList = append(videoList, video)
	}

	return &api.DouyinFavoriteListResponse{
		StatusCode: 0,
		VideoList:  videoList,
	}, nil
}
