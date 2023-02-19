package service

import (
	"douyin/biz/model/api"
	"reflect"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	type args struct {
		userID     uint64
		infoUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinUserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserInfo(tt.args.userID, tt.args.infoUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		req *api.DouyinUserLoginRequest
	}
	tests := []struct {
		name       string
		args       args
		wantUserID int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := Login(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("Login() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		req *api.DouyinUserRegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *api.DouyinUserRegisterResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Register(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
