package service

import (
	"bytes"
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/util"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

	// 封装video
	video := &db.Video{
		AuthorID:      uint64(userID),
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}

	return db.CreateVideo(video)
}

func GetPublishVideos(userID uint64) (*api.DouyinPublishListResponse, error) {
	res := new(api.DouyinPublishListResponse)
	videoList := make([]*api.Video, 0)

	videos, err := db.GetVideosByAuthorID(userID)
	if err != nil {
		return nil, err
	}

	// find author and pack data
	for i := 0; i < len(videos); i++ {
		u, err := db.SelectUserByID(uint(videos[i].AuthorID))
		if err != nil {
			return nil, err
		}

		video, err := pack.Videos(videos[i], u,
			db.IsFollow(userID, uint64(u.ID)), db.IsFavorite(userID, uint64(videos[i].ID)))
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}

	res.VideoList = videoList

	hlog.Info("pack over")

	return res, nil
}
