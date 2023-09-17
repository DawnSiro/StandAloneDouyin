package rdb

import (
	"douyin/dal/model"
	"fmt"
	"testing"
	"time"
)

func TestFollow(t *testing.T) {
	InitTest()
	type args struct {
		user   *model.FanUserRedisData
		toUser *model.FollowUserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "6-7", args: args{
			user: &model.FanUserRedisData{
				UID:      6,
				Username: "user01",
				Avatar:   "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			},
			toUser: &model.FollowUserData{
				UID:         7,
				Username:    "user02",
				Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
				CreatedTime: time.Now(),
			},
		}}, {name: "7-6", args: args{
			user: &model.FanUserRedisData{
				UID:      7,
				Username: "user02",
				Avatar:   "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			},
			toUser: &model.FollowUserData{
				UID:         6,
				Username:    "user01",
				Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
				CreatedTime: time.Now(),
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Follow(tt.args.user, tt.args.toUser); (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCancelFollow(t *testing.T) {
	InitTest()
	type args struct {
		user   *model.FanUserRedisData
		toUser *model.FollowUserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "user01", args: args{
			user: &model.FanUserRedisData{
				UID:      6,
				Username: "user01",
				Avatar:   "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			},
			toUser: &model.FollowUserData{
				UID:         7,
				Username:    "user02",
				Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
				CreatedTime: time.Now(),
			},
		}},
		{name: "user02", args: args{
			user: &model.FanUserRedisData{
				UID:      7,
				Username: "user02",
				Avatar:   "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
			},
			toUser: &model.FollowUserData{
				UID:         6,
				Username:    "user01",
				Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
				CreatedTime: time.Now(),
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CancelFollow(tt.args.user, tt.args.toUser); (err != nil) != tt.wantErr {
				t.Errorf("CancelFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetFollowZSet(t *testing.T) {
	InitTest()
	type args struct {
		userID     uint64
		followList []*model.FollowUserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "user01", args: args{
			userID: 6,
			followList: []*model.FollowUserData{
				{
					UID:         7,
					Username:    "user02",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				}, {
					UID:         8,
					Username:    "user03",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				},
			},
		}},
		{name: "user02", args: args{
			userID: 7,
			followList: []*model.FollowUserData{
				{
					UID:         6,
					Username:    "user01",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				}, {
					UID:         8,
					Username:    "user03",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFollowUserZSet(tt.args.userID, tt.args.followList); (err != nil) != tt.wantErr {
				t.Errorf("SetFollowUserZSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFollowUserZSet(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.FollowUserRedisData
		wantErr bool
	}{
		{name: "user01", args: args{userID: 6}},
		{name: "user02", args: args{userID: 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFollowUserZSet(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFollowUserZSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("------------")
			for _, data := range got {
				fmt.Println(data.UID)
				fmt.Println(data.Username)
				fmt.Println(data.Avatar)
				fmt.Println("------------")
			}
		})
	}
}

func TestSetFanZSet(t *testing.T) {
	InitTest()
	type args struct {
		userID     uint64
		followList []*model.FanUserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "user01", args: args{
			userID: 6,
			followList: []*model.FanUserData{
				{
					UID:         7,
					Username:    "user02",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				}, {
					UID:         8,
					Username:    "user03",
					Avatar:      "https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png",
					CreatedTime: time.Now(),
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFanUserZSet(tt.args.userID, tt.args.followList); (err != nil) != tt.wantErr {
				t.Errorf("SetFanUserZSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFanUserZSet(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.FanUserRedisData
		wantErr bool
	}{
		{name: "user01", args: args{userID: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFanUserZSet(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFanUserZSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("------------")
			for _, data := range got {
				fmt.Println(data.UID)
				fmt.Println(data.Username)
				fmt.Println(data.Avatar)
				fmt.Println("------------")
			}
		})
	}
}

func TestIsFollow(t *testing.T) {
	InitTest()
	type args struct {
		userID   uint64
		toUserID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "user03-user02", args: args{
			userID:   8,
			toUserID: 6,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsFollow(tt.args.userID, tt.args.toUserID)
			fmt.Println(tt.args, got, err)
		})
	}
}

func TestSetFollowUserIDSet(t *testing.T) {
	InitTest()
	type args struct {
		userID      uint64
		followIDSet []uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "user01", args: args{
			userID:      6,
			followIDSet: []uint64{7, 8},
		}},
		{name: "user02", args: args{
			userID:      7,
			followIDSet: []uint64{6, 8},
		}},
		{name: "user03", args: args{
			userID:      8,
			followIDSet: []uint64{6, 7},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFollowUserIDSet(tt.args.userID, tt.args.followIDSet); (err != nil) != tt.wantErr {
				t.Errorf("SetFollowUserIDSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFollowUserIDSet(t *testing.T) {
	InitTest()
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    map[uint64]struct{}
		wantErr bool
	}{
		{name: "user01", args: args{userID: 6}},
		{name: "user02", args: args{userID: 7}},
		{name: "user03", args: args{userID: 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFollowUserIDSet(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFollowUserIDSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("------------")
			for data := range got {
				fmt.Println(data)
			}
			fmt.Println("------------")
		})
	}
}
