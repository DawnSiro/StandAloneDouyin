package db

import (
	"douyin/pkg/initialize"
	"reflect"
	"testing"
)

func TestCancelFavoriteVideo(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID  uint64
		videoID uint64
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
			if err := CancelFavoriteVideo(tt.args.userID, tt.args.videoID); (err != nil) != tt.wantErr {
				t.Errorf("CancelFavoriteVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFavoriteVideo(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID  uint64
		videoID uint64
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
			if err := FavoriteVideo(tt.args.userID, tt.args.videoID); (err != nil) != tt.wantErr {
				t.Errorf("FavoriteVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsFavoriteVideo(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID  uint64
		videoID uint64
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
			if got := IsFavoriteVideo(tt.args.userID, tt.args.videoID); got != tt.want {
				t.Errorf("IsFavoriteVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectFavoriteVideoListByUserID(t *testing.T) {
	initialize.MySQL()
	type args struct {
		toUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*Video
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectFavoriteVideoListByUserID(tt.args.toUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectFavoriteVideoListByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectFavoriteVideoListByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserFavoriteVideo_TableName(t *testing.T) {
	initialize.MySQL()
	type fields struct {
		ID        uint64
		UserID    uint64
		VideoID   uint64
		IsDeleted uint8
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
			n := &UserFavoriteVideo{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				VideoID:   tt.fields.VideoID,
				IsDeleted: tt.fields.IsDeleted,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
