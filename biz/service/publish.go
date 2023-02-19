package service

import (
	"bytes"
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/util"
	"errors"
	"github.com/gofrs/uuid"
	"io"
)

func PublishAction(title string, videoData []byte, userID int64) error {
	if userID <= 0 {
		return errors.New("userID error")
	}

	// 上传 Object 需要一个实现了 io.Reader 接口的结构体
	var reader io.Reader = bytes.NewReader(videoData)
	u1, err := uuid.NewV4()
	if err != nil {
		return err
	}
	fileName := u1.String() + "." + "mp4"
	// 上传视频并生成封面
	playURL, coverURL, err := util.UploadVideo(&reader, fileName)
	if err != nil {
		return err
	}

	return db.CreateVideo(&db.Video{
		AuthorID:      uint64(userID),
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
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(videos[i].AuthorID)
		if err != nil {
			return nil, err
		}

		video := pack.Video(videos[i], u,
			db.IsFollow(userID, u.ID), db.IsFavoriteVideo(userID, videos[i].ID))
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}

	return &api.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		VideoList:  videoList,
	}, nil
}
