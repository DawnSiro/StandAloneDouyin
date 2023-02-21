package db

import (
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
}

func TestCreateUser(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name       string
		args       args
		wantUserID int64
		wantErr    bool
	}{
		{"Normal", args{user: &User{ID: 100, Username: "testUser01", Password: "123456"}}, 100, false},
		{"Duplicate ID", args{user: &User{ID: 100, Username: "testUser01", Password: "123456"}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := CreateUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("CreateUser() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestDecreaseUserFavoriteCount(t *testing.T) {
	type args struct {
		userID uint64
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
			got, err := DecreaseUserFavoriteCount(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecreaseUserFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecreaseUserFavoriteCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncreaseUserFavoriteCount(t *testing.T) {
	type args struct {
		userID uint64
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
			got, err := IncreaseUserFavoriteCount(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncreaseUserFavoriteCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncreaseUserFavoriteCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectUserByID(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectUserByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectUserByName(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    []*User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectUserByName(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectUserByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_TableName(t *testing.T) {
	type fields struct {
		ID             uint64
		Username       string
		Password       string
		FollowingCount uint64
		FollowerCount  uint64
		Avatar         string
		WorkCount      uint64
		FavoriteCount  uint64
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
			n := &User{
				ID:             tt.fields.ID,
				Username:       tt.fields.Username,
				Password:       tt.fields.Password,
				FollowingCount: tt.fields.FollowingCount,
				FollowerCount:  tt.fields.FollowerCount,
				Avatar:         tt.fields.Avatar,
				WorkCount:      tt.fields.WorkCount,
				FavoriteCount:  tt.fields.FavoriteCount,
			}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
