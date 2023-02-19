package db

import (
	"reflect"
	"testing"
	"time"
)

func TestComment_TableName(t *testing.T) {
	type fields struct {
		ID          uint64
		IsDeleted   uint8
		VideoID     uint64
		UserID      uint64
		Content     string
		CreatedTime time.Time
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
			n := &Comment{
				ID:          tt.fields.ID,
				IsDeleted:   tt.fields.IsDeleted,
				VideoID:     tt.fields.VideoID,
				UserID:      tt.fields.UserID,
				Content:     tt.fields.Content,
				CreatedTime: tt.fields.CreatedTime,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	type args struct {
		videoID uint64
		content string
		userID  uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateComment(tt.args.videoID, tt.args.content, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteCommentByID(t *testing.T) {
	type args struct {
		commentID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteCommentByID(tt.args.commentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCommentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteCommentByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCommentCreatedByMyself(t *testing.T) {
	type args struct {
		userID    uint64
		commentID uint64
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
			if got := IsCommentCreatedByMyself(tt.args.userID, tt.args.commentID); got != tt.want {
				t.Errorf("IsCommentCreatedByMyself() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectCommentListByVideoID(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*Comment
		wantErr bool
	}{
		{"Normal", args{videoID: 1}, []*Comment{
			&Comment{ID: 1, IsDeleted: 0, VideoID: 1, UserID: 1, Content: "Content01", CreatedTime: time.Now()}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectCommentListByVideoID(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCommentListByVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectCommentListByVideoID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
