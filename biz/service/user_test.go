package service

import (
	"douyin/biz/model/api"
	"douyin/pkg/initialize"
	"reflect"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID     uint64
		infoUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinUserResponse
		wantErr bool
	}{
		{"Normal",
			args{
				userID:     3,
				infoUserID: 3,
			}, &api.DouyinUserResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				User: &api.UserInfo{
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
			}, false},
		{"infoUserID_not_exist",
			args{
				userID:     3,
				infoUserID: 30000,
			}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserInfo(tt.args.userID, tt.args.infoUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	initialize.MySQL()
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name       string
		args       args
		wantUserID uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Login(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if uint64(resp.UserID) != tt.wantUserID {
				t.Errorf("Login() gotUserID = %v, want %v", resp, tt.wantUserID)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	initialize.MySQL()
	type args struct {
		username string
		password string
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
			if _, err := Register(tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
