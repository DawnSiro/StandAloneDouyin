package service

import (
	"context"
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/model"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/pulsar"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

// FavoriteVideo 点赞视频
func FavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {
	logTag := "service.favorite.FavoriteVideo err:"

	videoLikeCountKey := constant.FavoriteNumLimitPrefix + strconv.FormatUint(userID, 10)
	// 判断是否到了单位时间内的点赞上限
	userVideoLikeCount, err := global.VideoRC.Get(videoLikeCountKey).Result()
	var videoLikeCountInt int
	if err == redis.Nil {
		global.VideoRC.Set(videoLikeCountKey, "0", constant.VideoLikeLimitTime)
	} else {
		videoLikeCountInt, _ = strconv.Atoi(userVideoLikeCount)
		if videoLikeCountInt >= constant.VideoLikeLimit {
			return &api.DouyinFavoriteActionResponse{
				StatusCode: errno.VideoLikeLimitError.ErrCode,
				StatusMsg:  &errno.VideoLikeLimitError.ErrMsg,
			}, nil
		}
	}

	// 放入消息队列
	err = pulsar.GetLikeActionMQInstance().LikeAction(userID, videoID)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	hlog.Debug("service.favorite.FavoriteVideo: publish a message")

	// 更新单位时间内的点赞数量
	userVideoLikeCountTime, err := global.VideoRC.TTL(videoLikeCountKey).Result()
	global.VideoRC.Set(videoLikeCountKey, videoLikeCountInt+1, userVideoLikeCountTime)

	return &api.DouyinFavoriteActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

// CancelFavoriteVideo 取消点赞视频
func CancelFavoriteVideo(userID, videoID uint64) (*api.DouyinFavoriteActionResponse, error) {

	// 加入消息队列
	err := pulsar.GetLikeActionMQInstance().CancelLikeAction(userID, videoID)
	if err != nil {
		hlog.Error("service.favorite.CancelFavoriteVideo err:", err.Error())
		return nil, err
	}
	hlog.Debug("service.favorite.CancelFavoriteVideo: publish a message")
	return &api.DouyinFavoriteActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetFavoriteList(userID, selectUserID uint64) (*api.DouyinFavoriteListResponse, error) {
	logTag := "service.favorite.GetFavoriteList err:"

	//// 直接使用 JOIN 连接多个表数据，一次性查出所有数据。
	//videoList, err := db.SelectFavoriteVideoDataListByUserID(userID, selectUserID)
	//if err != nil {
	//	hlog.Error(logTag, err.Error())
	//	return nil, err
	//}

	// 加分布式锁
	lock := rdb.NewUserKeyLock(userID, constant.FavoriteVideoIDRedisZSetPrefix)
	// 如果 redis 不可用，应该使用程序代码之外的方式进行限流
	_ = lock.Lock(context.Background(), global.VideoRC)

	// 查询用户点赞视频ID列表
	ufVideoIDList, err := rdb.GetFavoriteVideoID(selectUserID)
	if err != nil {
		// 缓存未命中则查询数据库
		set, err := db.SelectFavoriteVideoIDZSet(selectUserID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
		rdbZSet := make([]*rdb.FavoriteVideoIDZSet, len(set))
		ufVideoIDList = make([]uint64, len(set))
		for i, id := range set {
			ufVideoIDList[i] = id.VideoID
			rdbZSet[i] = &rdb.FavoriteVideoIDZSet{
				VideoID:     id.VideoID,
				CreatedTime: float64(id.CreatedTime.UnixMilli()),
			}
		}
		// 设置缓存
		err = rdb.SetFavoriteVideoID(userID, rdbZSet)
		if err != nil {
			hlog.Error(logTag, err.Error())
		}
	}

	// 查询VideoInfo
	videoInfoList := make([]*model.Video, len(ufVideoIDList))
	lostVideoIDList := make([]uint64, 0)
	for i, u := range ufVideoIDList {
		info, err := rdb.GetVideoInfo(u)
		if err != nil {
			hlog.Error(logTag, err.Error())
			// 如果缓存不存在，则记录下
			lostVideoIDList = append(lostVideoIDList, u)
		}
		// info 为 nil 也先占着位置
		videoInfoList[i] = info
	}

	// 进行批处理查询
	lostVideoList, err := db.SelectVideoListByVideoID(lostVideoIDList)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	i := 0
	for _, v := range lostVideoList {
		for ; i < len(videoInfoList); i++ {
			// 找到空余的则填入
			if videoInfoList[i] == nil {
				videoInfoList[i] = v
				// 顺手设置缓存
				err := rdb.SetVideoInfo(v)
				if err != nil {
					hlog.Error(logTag, err.Error())
				}
			}
		}
	}

	// 查询UserInfo，需要注意可能有重复的视频作者
	userInfoList := make([]*model.User, len(videoInfoList))
	lostUserInfoIDList := make([]uint64, 0)
	for i, video := range videoInfoList {
		info, err := rdb.GetUserInfo(video.AuthorID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			// 如果缓存不存在，则记录下
			lostUserInfoIDList = append(lostUserInfoIDList, video.AuthorID)
		}
		// info 为 nil 也先占着位置
		userInfoList[i] = info
	}

	// 还是一次性查询完剩下的
	// 这里需要注意一个多个评论可能对应同一个用户，找坑的时候需要额外判断下
	lostUInfoList, err := db.SelectUserByIDList(lostUserInfoIDList)
	if err != nil {
		hlog.Error(logTag, err)
		return nil, err
	}

	i = 0
	for _, data := range lostUInfoList {
		for ; i < len(userInfoList); i++ {
			// 找到坑了则填入
			if userInfoList[i] == nil && videoInfoList[i].AuthorID == data.ID {
				userInfoList[i] = data
			}
		}
		// 顺手设置缓存
		err = rdb.SetUserInfo(data)
		if err != nil {
			hlog.Error(logTag, err)
		}
	}

	// 查询用户的点赞列表
	ufIDList, err := rdb.GetFavoriteVideoID(userID)
	ufIDSet := make(map[uint64]struct{}, len(ufIDList))
	for _, u := range ufIDList {
		ufIDSet[u] = struct{}{}
	}

	// 查询用户的关注列表
	followUserIDSet, err := rdb.GetFollowUserIDSet(userID)
	if err != nil {
		set, err := db.SelectFollowUserIDSet(userID)
		if err != nil {
			hlog.Error(logTag, err)
			return nil, err
		}
		err = rdb.SetFollowUserIDSet(userID, set)
		if err != nil {
			hlog.Error(logTag, err)
		}
	}

	// 解锁
	err = lock.Unlock(global.VideoRC)

	vList := make([]*api.Video, len(videoInfoList))
	for i := 0; i < len(videoInfoList); i++ {
		isFollow := false
		if _, ok := followUserIDSet[videoInfoList[i].AuthorID]; ok {
			isFollow = true
		}
		isFavorite := false
		if _, ok := ufIDSet[videoInfoList[i].AuthorID]; ok {
			isFavorite = true
		}
		vList[i] = pack.Video(videoInfoList[i], userInfoList[i], isFollow, isFavorite)
	}

	return &api.DouyinFavoriteListResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  vList,
	}, nil
}
