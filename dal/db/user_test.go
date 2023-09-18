package db

import (
	"douyin/dal/model"
	"testing"
)

func TestCreateUser(t *testing.T) {
	InitTest()
	type args struct {
		user *model.User
	}
	tests := []struct {
		name       string
		args       args
		wantUserID int64
		wantErr    bool
	}{
		{"Normal", args{user: &model.User{ID: 100, Username: "testUser01", Password: "123456"}}, 100, false},
		{"Duplicate ID", args{user: &model.User{ID: 100, Username: "testUser01", Password: "123456"}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := CreateUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if int64(gotUserID) != tt.wantUserID {
				t.Errorf("CreateUser() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
