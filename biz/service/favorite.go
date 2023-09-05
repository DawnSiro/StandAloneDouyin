package service

import (
	"douyin/dal/pack"
	"douyin/pkg/pulsar"
	"strconv"
	"strings"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

// FavoriteVideo this is a func for add Favorite or reduce Favorite
func FavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_like")
	videoLikeKey := builder.String()

	var builder1 strings.Builder
	builder1.WriteString(strconv.FormatUint(userID, 10))
	builder1.WriteString("_user_like_count")
	videoLikeCountKey := builder1.String()

	userVideoLikeCount, err := global.VideoFRC.Get(videoLikeCountKey).Result()
	var videoLikeCountInt int
	if err == redis.Nil {
		global.VideoFRC.Set(videoLikeCountKey, "0", constant.VideoLikeLimitTime)
	} else {
		videoLikeCountInt, _ = strconv.Atoi(userVideoLikeCount)
		if videoLikeCountInt >= constant.VideoLikeLimit {
			return &api.DouyinFavoriteActionResponse{
				StatusCode: errno.VideoLikeLimitError.ErrCode,
				StatusMsg:  &errno.VideoLikeLimitError.ErrMsg,
			}, nil
		}
	}

	likeCount, err := global.VideoFRC.Get(videoLikeKey).Result()
	if err == redis.Nil {
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
			return nil, err
		}
		global.VideoFRC.Set(videoLikeKey, likeInt64, 0)
	}
	var likeUint64 uint64
	if likeCount != "" {
		likeUint64, err = strconv.ParseUint(likeCount, 10, 64)
		if err != nil {
			hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
			return nil, err
		}
	}

	err = pulsar.GetLikeActionMQInstance().LikeAction(userID, videoID)
	if err != nil {
		hlog.Error("service.favorite.FavoriteVideo err:", err.Error())
		return nil, err
	}
	hlog.Debug("service.favorite.FavoriteVideo: publish a message")

	// 如果 DB 层事务回滚了，err 就不为 nil，Redis 里的数据就不会更新
	global.VideoFRC.Set(videoLikeKey, likeUint64+1, 0)
	// 更新单位时间内的点赞数量
	userVideoLikeCountTime, err := global.VideoFRC.TTL(videoLikeCountKey).Result()
	global.VideoFRC.Set(videoLikeCountKey, videoLikeCountInt+1, userVideoLikeCountTime)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func CancelFavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	var builder strings.Builder
	builder.WriteString(strconv.FormatUint(videoID, 10))
	builder.WriteString("_video_like")
	videoLikeKey := builder.String()

	likeCount, err := global.VideoFRC.Get(videoLikeKey).Result()
	if err == redis.Nil {
		likeInt64, err := db.SelectVideoFavoriteCountByVideoID(videoID)
		if err != nil {
			hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
			return nil, err
		}
		global.VideoFRC.Set(videoLikeKey, likeInt64, 0)
	}

	var likeUint64 uint64
	if likeCount != "" {
		likeUint64, err = strconv.ParseUint(likeCount, 10, 64)
		if err != nil {
			hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
			return nil, err
		}
	}

	err = pulsar.GetLikeActionMQInstance().CancelLikeAction(userID, videoID)
	if err != nil {
		hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
		return nil, err
	}
	hlog.Debug("service.favorite.CancelFavoriteVideo: publish a message")
	// 如果 DB 层事务回滚了，err 就不为 nil，Redis 里的数据就不会更新
	global.VideoFRC.Set(videoLikeKey, likeUint64-1, 0)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetFavoriteList(userID, selectUserID uint64) (*api.DouyinFavoriteListResponse, error) {
	// 方案一，每次都按单表查询，存在循环查询数据的问题，经过测试，开启协程进行异步也没有多少性能提升
	//videos, err := db.SelectFavoriteVideoListByUserID(selectUserID)
	//if err != nil {
	//	hlog.Error("service.favorite.GetFavoriteList err:", err.Error())
	//	return nil, err
	//}
	//
	//// TODO 优化循环查询数据库问题
	//videoList := make([]*api.Video, 0)
	//for i := 0; i < len(videos); i++ {
	//	u, err := db.SelectUserByID(videos[i].AuthorID)
	//	if err != nil {
	//		hlog.Error("service.favorite.GetFavoriteList err:", err.Error())
	//		return nil, err
	//	}
	//	video := pack.Video(videos[i], u,
	//		db.IsFollow(userID, selectUserID), db.IsFavoriteVideo(userID, videos[i].ID))
	//	videoList = append(videoList, video)
	//}

	// 方案二，直接使用 JOIN 连接多个表数据，一次性查出所有数据。
	videoList, err := db.SelectFavoriteVideoDataListByUserID(userID, selectUserID)
	if err != nil {
		return nil, err
	}

	return &api.DouyinFavoriteListResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  pack.VideoDataList(videoList),
	}, nil
}
