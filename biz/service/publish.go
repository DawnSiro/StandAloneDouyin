package service

import (
	"bytes"
	"errors"
	"io"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gofrs/uuid"
)

func PublishAction(title string, videoData []byte, userID uint64) error {
	if userID == 0 {
		err := errors.New("userID error")
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return err
	}

	// 上传 Object 需要一个实现了 io.Reader 接口的结构体
	var reader io.Reader = bytes.NewReader(videoData)
	u1, err := uuid.NewV4()
	if err != nil {
		hlog.Error("service.publish.PublishAction err:", err.Error())
		return err
	}
	fileName := u1.String() + "." + "mp4"
	hlog.Info("service.publish.PublishAction videoName:", fileName)
	// 上传视频并生成封面
	playURL, coverURL, err := util.UploadVideo(&reader, fileName)
	if err != nil {
		return err
	}

	return db.CreateVideo(&db.Video{
		AuthorID:      userID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
}

func GetPublishVideos(userID uint64) (*api.DouyinPublishListResponse, error) {
	videoList := make([]*api.Video, 0)

	videos, err := db.GetVideosByAuthorID(userID)
	if err != nil {
		hlog.Error("service.publish.GetPublishVideos err:", err.Error())
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			hlog.Error("service.publish.GetPublishVideos err:", err.Error())
			return nil, err
		}

		video := pack.Video(videos[i], u,
			db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
		videoList = append(videoList, video)
	}

	return &api.DouyinPublishListResponse{
		StatusCode: 0,
		VideoList:  videoList,
	}, nil
}
