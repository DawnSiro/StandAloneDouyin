package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestGetMessageChat(t *testing.T) {
	type args struct {
		userID     uint64
		oppositeID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinMessageChatResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessageChat(tt.args.userID, tt.args.oppositeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessageChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMessageChat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	type args struct {
		fromUserID uint64
		toUserID   uint64
		content    string
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinMessageActionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendMessage(tt.args.fromUserID, tt.args.toUserID, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
