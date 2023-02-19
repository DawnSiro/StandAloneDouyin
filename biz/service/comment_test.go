package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestCommentList(t *testing.T) {
	type args struct {
		userID  uint64
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinCommentListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CommentList(tt.args.userID, tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	type args struct {
		userID    uint64
		videoID   uint64
		commentID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinCommentActionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteComment(tt.args.userID, tt.args.videoID, tt.args.commentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostComment(t *testing.T) {
	type args struct {
		userID      uint64
		videoID     uint64
		commentText string
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinCommentActionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostComment(tt.args.userID, tt.args.videoID, tt.args.commentText)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
