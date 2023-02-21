package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"reflect"
	"testing"
)

func TestCancelFavoriteVideo(t *testing.T) {
	type args struct {
		userID  uint64
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinFavoriteActionResponse
		wantErr bool
	}{
		// TODO: 这里取消点赞不行
		{"Normal", args{
			userID:  101,
			videoID: 6,
		}, &api.DouyinFavoriteActionResponse{
			StatusCode: 0,
			StatusMsg:  nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CancelFavoriteVideo(tt.args.userID, tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelFavoriteVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CancelFavoriteVideo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFavoriteList(t *testing.T) {
	db.Init()
	v1 := int64(0)
	type args struct {
		userID       uint64
		selectUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinFavoriteListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Normal",
			args{
				userID:       101,
				selectUserID: 101,
			},
			&api.DouyinFavoriteListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				VideoList: []*api.Video{
					{
						ID: 6,
						Author: &api.UserInfo{
							ID:            3,
							Name:          "user01",
							FollowCount:   v1,
							FollowerCount: v1,
							IsFollow:      true,
							Avatar:        "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
						},
						PlayURL:       "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/54709da0-71bf-488e-ab89-fb11db5ff1ae.mp4",
						CoverURL:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/",
						FavoriteCount: 2,
						CommentCount:  5,
						IsFavorite:    true,
						Title:         "video03",
					},
				},
			},
			false},
		{"userID_err",
			args{
				userID:       101,
				selectUserID: 1000000,
			}, &api.DouyinFavoriteListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				VideoList:  []*api.Video{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FavoriteList(tt.args.userID, tt.args.selectUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FavoriteList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFavoriteVideo(t *testing.T) {
	type args struct {
		userID  uint64
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinFavoriteActionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		// TODO: 这里正常点赞不行
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FavoriteVideo(tt.args.userID, tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavoriteVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FavoriteVideo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
