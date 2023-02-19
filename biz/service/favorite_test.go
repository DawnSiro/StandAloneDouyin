package service

import (
	"douyin/biz/model/api"
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
		// TODO: Add test cases.
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
