package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestGetFeed(t *testing.T) {
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
		// TODO: Add test cases.
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
