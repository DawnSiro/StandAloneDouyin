package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"reflect"
	"testing"
)

func TestCommentList(t *testing.T) {
	db.Init()
	type args struct {
		userID  uint64
		videoID uint64
	}
	v1 := int64(0)
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinCommentListResponse
		wantErr bool
	}{
		{"Normal", args{
			userID:  101,
			videoID: 6,
		}, &api.DouyinCommentListResponse{StatusCode: 0, StatusMsg: nil, CommentList: []*api.Comment{{
			ID: 39,
			User: &api.User{
				ID:            101,
				Name:          "testUser1",
				FollowCount:   &v1,
				FollowerCount: &v1,
				IsFollow:      true,
				Avatar:        "",
			},
			Content:    "测试啦",
			CreateDate: "02-20",
		}}}, false},
		{"videoId_err", args{
			userID:  101,
			videoID: 10000,
		}, &api.DouyinCommentListResponse{StatusCode: 0, StatusMsg: nil, CommentList: []*api.Comment{}}, false},
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
	db.Init()
	v1 := int64(0)
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
		{"Normal", args{
			userID:    101,
			videoID:   6,
			commentID: 39,
		}, &api.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: nil, Comment: &api.Comment{
			ID: 39,
			User: &api.User{
				ID:            101,
				Name:          "ceshi1",
				FollowCount:   &v1,
				FollowerCount: &v1,
				IsFollow:      false,
				Avatar:        "",
			},
			Content:    "测试啦",
			CreateDate: "02-20",
		}}, false},
		//TODO: 这里出现问题了！
		{"video_id_err", args{
			userID:    101,
			videoID:   10000,
			commentID: 38,
		}, &api.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: nil, Comment: nil}, true},
		{"comment_id_err", args{
			userID:    101,
			videoID:   6,
			commentID: 10000,
		}, &api.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: nil, Comment: nil}, true},
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
	db.Init()
	v1 := int64(0)
	//v2 := "评论不能为空"
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
		{"Normal", args{
			userID:      101,
			videoID:     6,
			commentText: "测试啦",
		}, &api.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: nil, Comment: &api.Comment{
			ID: 43,
			User: &api.User{
				ID:            101,
				Name:          "ceshi1",
				FollowCount:   &v1,
				FollowerCount: &v1,
				IsFollow:      false,
				Avatar:        "",
			},
			Content:    "测试啦",
			CreateDate: "02-20",
		}}, false},
		{"commentText_err", args{
			userID:      101,
			videoID:     6,
			commentText: "",
		}, &api.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: nil, Comment: nil}, true},
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
