package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/model"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gofrs/uuid"
)

func PublishAction(title string, videoData []byte, userID uint64) (*api.DouyinPublishActionResponse, error) {
	logTag := "service.publish.PublishAction err:"
	if userID == 0 {
		err := errors.New("userID error")
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	// 上传 Object 需要一个实现了 io.Reader 接口的结构体
	var reader io.Reader = bytes.NewReader(videoData)
	u1, err := uuid.NewV4()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"
	hlog.Info("service.publish.PublishAction videoName:", fileName)
	// 上传视频并生成封面
	playURL, coverURL, err := util.UploadVideo(&reader, fileName)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	videoID, err := db.CreateVideo(&model.Video{
		PublishTime:   time.Now(),
		AuthorID:      userID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	// 加入布隆过滤器
	global.VideoIDBloomFilter.AddString(strconv.FormatUint(videoID, 10))

	return &api.DouyinPublishActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetPublishVideos(userID, selectUserID uint64) (*api.DouyinPublishListResponse, error) {
	logTag := "service.favorite.GetPublishVideos err:"

	//videoData, err := db.SelectPublishVideoDataListByUserID(userID, selectUserID)
	//if err != nil {
	//	hlog.Error("service.publish.GetPublishVideos err:", err.Error())
	//	return nil, err
	//}

	// 加分布式锁
	lock := rdb.NewUserKeyLock(userID, constant.FavoriteVideoIDRedisZSetPrefix)
	// 如果 redis 不可用，应该使用程序代码之外的方式进行限流
	_ = lock.Lock(context.Background(), global.VideoRC)

	// 查询用户点赞视频ID列表
	ufVideoIDList, err := rdb.GetPublishVideoID(selectUserID)
	if err != nil {
		// 缓存未命中则查询数据库
		set, err := db.SelectPublishVideoIDZSet(selectUserID)
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
		rdbZSet := make([]*rdb.PublishVideoIDZSet, len(set))
		ufVideoIDList = make([]uint64, len(set))
		for i, id := range set {
			ufVideoIDList[i] = id.VideoID
			rdbZSet[i] = &rdb.PublishVideoIDZSet{
				VideoID:     id.VideoID,
				CreatedTime: float64(id.CreatedTime.UnixMilli()),
			}
		}
		// 设置缓存
		err = rdb.SetPublishVideoID(userID, rdbZSet)
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

	return &api.DouyinPublishListResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  vList,
	}, nil
}
