package db

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateVideo(t *testing.T) {
	type args struct {
		video *Video
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
			if err := CreateVideo(tt.args.video); (err != nil) != tt.wantErr {
				t.Errorf("CreateVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecreaseCommentCount(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecreaseCommentCount(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecreaseCommentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecreaseCommentCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecreaseVideoFavoriteCount(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecreaseVideoFavoriteCount(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecreaseVideoFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecreaseVideoFavoriteCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVideosByAuthorID(t *testing.T) {
	type args struct {
		userID uint64
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
			got, err := GetVideosByAuthorID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideosByAuthorID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVideosByAuthorID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncreaseCommentCount(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncreaseCommentCount(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncreaseCommentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncreaseCommentCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncreaseVideoFavoriteCount(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncreaseVideoFavoriteCount(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncreaseVideoFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncreaseVideoFavoriteCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMGetVideos(t *testing.T) {
	type args struct {
		maxVideoNum int
		latestTime  *int64
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
			got, err := MGetVideos(tt.args.maxVideoNum, tt.args.latestTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGetVideos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectAuthorIDByVideoID(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectAuthorIDByVideoID(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectAuthorIDByVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectAuthorIDByVideoID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectCommentCountByVideoID(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectCommentCountByVideoID(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCommentCountByVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectCommentCountByVideoID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectVideoFavoriteCountByVideoID(t *testing.T) {
	type args struct {
		videoID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectVideoFavoriteCountByVideoID(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectVideoFavoriteCountByVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectVideoFavoriteCountByVideoID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCommentCount(t *testing.T) {
	type args struct {
		videoID      uint64
		commentCount uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateCommentCount(tt.args.videoID, tt.args.commentCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCommentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateCommentCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateVideoFavoriteCount(t *testing.T) {
	type args struct {
		videoID       uint64
		favoriteCount uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateVideoFavoriteCount(tt.args.videoID, tt.args.favoriteCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideoFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateVideoFavoriteCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_TableName(t *testing.T) {
	type fields struct {
		ID            uint64
		PublishTime   time.Time
		AuthorID      uint64
		PlayURL       string
		CoverURL      string
		FavoriteCount uint64
		CommentCount  uint64
		Title         string
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
			n := &Video{
				ID:            tt.fields.ID,
				PublishTime:   tt.fields.PublishTime,
				AuthorID:      tt.fields.AuthorID,
				PlayURL:       tt.fields.PlayURL,
				CoverURL:      tt.fields.CoverURL,
				FavoriteCount: tt.fields.FavoriteCount,
				CommentCount:  tt.fields.CommentCount,
				Title:         tt.fields.Title,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
