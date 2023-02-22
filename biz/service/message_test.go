package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"reflect"
	"testing"
)

func TestGetMessageChat(t *testing.T) {
	db.Init()
	v1 := int64(1676972817000)
	v2 := int64(1676972820000)
	type args struct {
		userID     uint64
		oppositeID uint64
		preMsgTime int64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinMessageChatResponse
		wantErr bool
	}{
		{"Normal", args{
			userID:     3,
			oppositeID: 5,
			preMsgTime: 0,
		},
			&api.DouyinMessageChatResponse{
				StatusCode: 0,
				StatusMsg:  nil,
				MessageList: []*api.Message{
					{
						ID:         1,
						ToUserID:   3,
						FromUserID: 5,
						Content:    "我是yuleng，你好",
						CreateTime: &v1,
					},
					{
						ID:         2,
						ToUserID:   5,
						FromUserID: 3,
						Content:    "我是user01，你好",
						CreateTime: &v2,
					},
				},
			}, false},
		{"toUserID_not_exist",
			args{
				userID:     300000,
				oppositeID: 5,
				preMsgTime: 0,
			},
			&api.DouyinMessageChatResponse{
				StatusCode:  0,
				StatusMsg:   nil,
				MessageList: []*api.Message{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMessageChat(tt.args.userID, tt.args.oppositeID, tt.args.preMsgTime)
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
