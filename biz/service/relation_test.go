package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestCancelFollow(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFollowListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
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
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFollowerListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
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
	type args struct {
		req *api.DouyinRelationFriendListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinRelationFriendListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFriendList(tt.args.req)
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
