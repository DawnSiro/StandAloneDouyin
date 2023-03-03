package service

import (
	"douyin/biz/model/api"
	"douyin/pkg/initialize"
	"reflect"
	"testing"
)

func TestGetFeed(t *testing.T) {
	initialize.MySQL()
	v1 := int64(40)
	v2 := int64(1676991666000)
	v3 := int64(1576991666000)
	type args struct {
		latestTime *int64
		userID     uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinFeedResponse
		wantErr bool
	}{
		{"Normal", args{
			latestTime: nil,
			userID:     3,
		}, &api.DouyinFeedResponse{
			StatusCode: 0,
			StatusMsg:  nil,
			VideoList: []*api.Video{
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
			},
			NextTime: &v1,
		}, false},
		{"latestTime_exist", args{
			latestTime: &v2,
			userID:     3,
		}, &api.DouyinFeedResponse{
			StatusCode: 0,
			StatusMsg:  nil,
			VideoList: []*api.Video{
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
			},
			NextTime: &v1,
		}, false},
		{"latestTime_so_early", args{
			latestTime: &v3,
			userID:     3,
		}, &api.DouyinFeedResponse{
			StatusCode: 0,
			StatusMsg:  nil,
			VideoList:  []*api.Video{},
			NextTime:   nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFeed(tt.args.latestTime, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFeed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
