package db

import (
	"douyin/dal/model"
	"fmt"
	"testing"
)

func TestSelectFriendDataList(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.FriendUserData
		wantErr bool
	}{
		{name: "user01", args: args{userID: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectFriendDataList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectFriendDataList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("------------")
			for _, data := range got {
				fmt.Println(data)
				fmt.Println("------------")
			}
		})
	}
}
