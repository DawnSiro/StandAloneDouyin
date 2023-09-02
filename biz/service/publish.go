package service

import (
	"bytes"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"io"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gofrs/uuid"
)

func PublishAction(title string, videoData []byte, userID uint64) (*api.DouyinPublishActionResponse, error) {
	if userID == 0 {
		err := errors.New("userID error")
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return nil, err
	}

	// 上传 Object 需要一个实现了 io.Reader 接口的结构体
	var reader io.Reader = bytes.NewReader(videoData)
	u1, err := uuid.NewV4()
	if err != nil {
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"
	hlog.Info("service.publish.PublishAction videoName:", fileName)
	// 上传视频并生成封面
	playURL, coverURL, err := util.UploadVideo(&reader, fileName)
	if err != nil {
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return nil, err
	}

	err = db.CreateVideo(&db.Video{
		PublishTime:   time.Now(),
		AuthorID:      userID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
	if err != nil {
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return nil, err
	}

	return &api.DouyinPublishActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetPublishVideos(userID, selectUserID uint64) (*api.DouyinPublishListResponse, error) {
	// 方案一，每次都按单表查询，存在循环查询数据的问题，经过测试，开启协程进行异步也没有多少性能提升
	//videoList := make([]*api.Video, 0)
	//
	//videos, err := db.GetVideosByAuthorID(userID)
	//if err != nil {
	//	hlog.Error("service.publish.GetPublishVideos err:", err.Error())
	//	return nil, err
	//}
	//
	//for i := 0; i < len(videos); i++ {
	//	u, err := db.SelectUserByID(videos[i].AuthorID)
	//	if err != nil {
	//		hlog.Error("service.publish.GetPublishVideos err:", err.Error())
	//		return nil, err
	//	}
	//
	//	video := pack.Video(videos[i], u,
	//		db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
	//	videoList = append(videoList, video)
	//}

	// Check if user info is available in Redis cache
	cacheKey := fmt.Sprintf("publishedVideoList:%d", selectUserID)
	cachedData, err := global.UserInfoRC.Get(cacheKey).Result()

	if err == nil {
		// Cache hit, judge if they are friends
		isFriend, err := AreUsersFriends(userID, selectUserID)
		if err != nil {
			hlog.Error("service.publish.GetPublishVideos err: Error checking if users are friends, ", err.Error())
		} else if isFriend {
			// Return cached published video list
			var cachedResponse api.DouyinPublishListResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
				hlog.Error("service.publish.GetPublishVideos err: Error decoding cached data, ", err.Error())
			} else {
				return &cachedResponse, nil
			}
		}
	}

	// Cache miss, query the database
	videoData, err := db.SelectPublishVideoDataListByUserID(userID, selectUserID)
	if err != nil {
		hlog.Error("service.publish.GetPublishVideos err:", err.Error())
		return nil, err
	}

	// Pack video data
	response := &api.DouyinPublishListResponse{
		StatusCode: errno.Success.ErrCode,
		VideoList:  pack.VideoDataList(videoData),
	}

	// Store the published video list in Redis cache
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, 24*time.Hour).Err()
	if err != nil {
		hlog.Error("service.publish.GetPublishVideos err: Error storing data in cache, ", err.Error())
	}

	return response, nil
}

func AreUsersFriends(userID, selectUserID uint64) (bool, error) {
	friendListResponse, err := GetFriendList(userID)
	if err != nil {
		return false, err
	}

	for _, friend := range friendListResponse.UserList {
		if friend.ID == int64(selectUserID) {
			return true, nil
		}
	}
	// Not friends
	return false, nil
}
