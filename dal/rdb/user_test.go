package rdb

import (
	"douyin/dal/model"
	"fmt"
	"testing"
)

func TestExpireUserInfo(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"user01", args{userID: 6}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExpireUserInfo(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ExpireUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"user01", args{userID: 6}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserInfo(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestSetUserInfo(t *testing.T) {
	InitTest()
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "user01", args: args{&model.User{
			ID:              6,
			Username:        "user01",
			FollowingCount:  0,
			FollowerCount:   0,
			Avatar:          "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			BackgroundImage: "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/background/background.png",
			Signature:       "",
			TotalFavorited:  0,
			WorkCount:       0,
			FavoriteCount:   0,
		}}},
		{name: "user02", args: args{&model.User{
			ID:              7,
			Username:        "user02",
			FollowingCount:  0,
			FollowerCount:   0,
			Avatar:          "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			BackgroundImage: "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/background/background.png",
			Signature:       "",
			TotalFavorited:  0,
			WorkCount:       0,
			FavoriteCount:   0,
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetUserInfo(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
