package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestGetPublishVideos(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinPublishListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
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
		userID    int64
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
