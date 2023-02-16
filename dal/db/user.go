package db

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"column:username;index:idx_username,unique;type:varchar(40);not null" json:"username"`
	Password       string `gorm:"type:varchar(256);not null" json:"password"`
	FollowingCount int64  `gorm:"default:0" json:"following_count"`
	FollowerCount  int64  `gorm:"default:0" json:"follower_count"`
	Avatar         string `gorm:"type:varchar(256)" json:"avatar"`
}

func (n *User) TableName() string {
	return constant.UserTableName
}

// CreateUser create user
func CreateUser(user *User) (userID int64, err error) {
	if err := DB.Create(user).Error; err != nil {
		return 0, err
	}
	return int64(user.ID), nil
}

func SelectUserByUserID(userID uint) (*api.User, error) {
	var result api.User

	user := &User{
		Model: gorm.Model{
			ID: userID,
		},
	}
	if err := DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	result.ID = int64(user.ID)
	result.Name = user.Username
	result.FollowCount = &user.FollowingCount
	result.FollowerCount = &user.FollowerCount
	result.IsFollow = false
	result.Avatar = user.Avatar

	return &result, nil
}

func SelectUserByID(userId uint) (*User, error) {
	res := &User{
		Model: gorm.Model{
			ID: userId,
		},
	}

	if err := DB.First(&res, "id = ?", userId).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func SelectUserByName(username string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.Where("username = ?", username).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
