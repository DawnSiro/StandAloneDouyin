package db

import (
	"douyin/biz/model/api"
	"douyin/constants"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName       string `json:"user_name"`
	Password       string `json:"password"`
	FollowingCount int64  `json:"following_count"`
	FollowerCount  int64  `json:"follower_count"`
}

func (n *User) TableName() string {
	return constants.UserTableName
}

func SelectUserByUserId(userId uint) (*api.User, error) {
	var result api.User

	user := &User{
		Model: gorm.Model{
			ID: userId,
		},
	}
	if err := DB.First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	result.ID = int64(user.ID)
	result.Name = user.UserName
	result.FollowCount = &user.FollowingCount
	result.FollowerCount = &user.FollowerCount
	result.IsFollow = false
	//TODO:miss avatar

	return &result, nil
}
