package db

import (
	"douyin/biz/handler/api"
	"douyin/pkg/initialize"
	"reflect"
	"testing"
	"time"
)

func TestCreateMessage(t *testing.T) {
	initialize.MySQL()
	type args struct {
		fromUserID uint64
		toUserID   uint64
		content    string
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
			if err := CreateMessage(tt.args.fromUserID, tt.args.toUserID, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLatestMsg(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID     uint64
		oppositeID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *FriendMessageResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestMsg(tt.args.userID, tt.args.oppositeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestMsg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMessagesByUserIDAndPreMsgTime(t *testing.T) {
	initialize.MySQL()
	type args struct {
		userID     uint64
		oppositeID uint64
		preMsgTime int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*api.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessagesByUserIDAndPreMsgTime(tt.args.userID, tt.args.oppositeID, tt.args.preMsgTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessagesByUserIDAndPreMsgTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMessagesByUserIDAndPreMsgTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_TableName(t *testing.T) {
	initialize.MySQL()
	type fields struct {
		ID         uint64
		ToUserID   uint64
		FromUserID uint64
		Content    string
		CreateTime time.Time
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
			n := &api.Message{
				ID:         tt.fields.ID,
				ToUserID:   tt.fields.ToUserID,
				FromUserID: tt.fields.FromUserID,
				Content:    tt.fields.Content,
				CreateTime: tt.fields.CreateTime,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
