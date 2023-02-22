package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"reflect"
	"testing"
)

func TestGetPublishVideos(t *testing.T) {
	db.Init()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinPublishListResponse
		wantErr bool
	}{
		{"Normal", args{userID: 3},
			&api.DouyinPublishListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				VideoList: []*api.Video{
					{
						ID: 1,
						Author: &api.UserInfo{
							ID:              3,
							Name:            "user01",
							FollowCount:     2,
							FollowerCount:   1,
							IsFollow:        true,
							Avatar:          "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
							BackgroundImage: "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
							Signature:       "",
							TotalFavorited:  0,
							WorkCount:       6,
							FavoriteCount:   1,
						},
						PlayURL:       "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/video01.mp4",
						CoverURL:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/cover01.png",
						FavoriteCount: 0,
						CommentCount:  3,
						IsFavorite:    false,
						Title:         "title01",
					},
					{
						ID: 2,
						Author: &api.UserInfo{
							ID:              3,
							Name:            "user01",
							FollowCount:     2,
							FollowerCount:   1,
							IsFollow:        true,
							Avatar:          "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
							BackgroundImage: "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
							Signature:       "",
							TotalFavorited:  0,
							WorkCount:       6,
							FavoriteCount:   1,
						},
						PlayURL:       "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/video02.mp4",
						CoverURL:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/cover02.png",
						FavoriteCount: 1,
						CommentCount:  17,
						IsFavorite:    true,
						Title:         "title02",
					},
				},
			}, false},
		{"userID_not_exist", args{userID: 300000},
			&api.DouyinPublishListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				VideoList:  []*api.Video{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublishVideos(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublishVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPublishVideos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublishAction(t *testing.T) {
	type args struct {
		title     string
		videoData []byte
		userID    uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PublishAction(tt.args.title, tt.args.videoData, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("PublishAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
