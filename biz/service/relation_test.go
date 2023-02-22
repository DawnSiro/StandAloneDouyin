package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"reflect"
	"testing"
)

func TestCancelFollow(t *testing.T) {
	db.Init()
	type args struct {
		userID   uint64
		toUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationActionResponse
		wantErr bool
	}{
		{"Normal",
			args{
				userID:   3,
				toUserID: 5,
			},
			&api.DouyinRelationActionResponse{
				StatusCode: 0,
				StatusMsg:  nil,
			},
			false},
		{"toUserID_not_exist",
			args{
				userID:   3,
				toUserID: 100000,
			}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CancelFollow(tt.args.userID, tt.args.toUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelFollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CancelFollow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFollow(t *testing.T) {
	db.Init()
	type args struct {
		userID   uint64
		toUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationActionResponse
		wantErr bool
	}{
		{"Normal",
			args{
				userID:   3,
				toUserID: 5,
			},
			&api.DouyinRelationActionResponse{
				StatusCode: 0,
				StatusMsg:  nil,
			}, false},
		{"toUserID_not_exist",
			args{
				userID:   3,
				toUserID: 100000,
			}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Follow(tt.args.userID, tt.args.toUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Follow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFollowList(t *testing.T) {
	db.Init()
	v1 := int64(2)
	v2 := int64(1)
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFollowListResponse
		wantErr bool
	}{
		{"Normal",
			args{userID: 5},
			&api.DouyinRelationFollowListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList: []*api.User{
					{
						ID:            3,
						Name:          "user01",
						FollowCount:   &v1,
						FollowerCount: &v2,
						IsFollow:      true,
						Avatar:        "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
					},
				},
			}, false},
		{"userId_not_exist",
			args{userID: 50000},
			&api.DouyinRelationFollowListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList:   []*api.User{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFollowList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFollowList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFollowList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFollowerList(t *testing.T) {
	db.Init()
	v1 := int64(2)
	v2 := int64(1)
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFollowerListResponse
		wantErr bool
	}{
		{"Normal",
			args{userID: 5},
			&api.DouyinRelationFollowerListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList: []*api.User{
					{
						ID:            3,
						Name:          "user01",
						FollowCount:   &v1,
						FollowerCount: &v2,
						IsFollow:      true,
						Avatar:        "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
					},
				},
			}, false},
		{"userID_not_exist",
			args{userID: 50000},
			&api.DouyinRelationFollowerListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList:   []*api.User{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFollowerList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFollowerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFollowerList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFriendList(t *testing.T) {
	db.Init()
	v1 := int64(2)
	v2 := int64(1)
	v3 := "你好"
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFriendListResponse
		wantErr bool
	}{
		{"Normal", args{userID: 5},
			&api.DouyinRelationFriendListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList: []*api.FriendUser{
					{
						ID:            3,
						Name:          "user01",
						FollowCount:   &v1,
						FollowerCount: &v2,
						IsFollow:      true,
						Avatar:        "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/%E6%B5%8B%E8%AF%95%E5%9B%BE%E7%89%871.png",
						Message:       &v3,
						MsgType:       0,
					},
				},
			}, false},
		{"userID_not_exist",
			args{userID: 500000}, &api.DouyinRelationFriendListResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				UserList:   []*api.FriendUser{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFriendList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFriendList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFriendList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
