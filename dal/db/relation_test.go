package db

import (
	"douyin/biz/model/api"
	"douyin/pkg/initialize"
	"reflect"
	"testing"
)

func TestCancelFollow(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID   uint64
		toUserID uint64
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
			if err := CancelFollow(tt.args.userID, tt.args.toUserID); (err != nil) != tt.wantErr {
				t.Errorf("CancelFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFollow(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID   uint64
		toUserID uint64
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
			if err := Follow(tt.args.userID, tt.args.toUserID); (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFollowList(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*User
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
	initialize.MySQL()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*User
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
	initialize.MySQL()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*api.FriendUser
		wantErr bool
	}{
		// TODO: Add test cases.
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

func TestIsFollow(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID   uint64
		toUserID uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFollow(tt.args.userID, tt.args.toUserID); got != tt.want {
				t.Errorf("IsFollow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_TableName(t *testing.T) {
	initialize.MySQL()
	type fields struct {
		ID        uint64
		IsDeleted uint8
		UserID    uint64
		ToUserID  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Relation{
				ID:        tt.fields.ID,
				IsDeleted: tt.fields.IsDeleted,
				UserID:    tt.fields.UserID,
				ToUserID:  tt.fields.ToUserID,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
