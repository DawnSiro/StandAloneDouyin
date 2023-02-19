package db

import (
	"douyin/pkg/constant"
)

type User struct {
	ID             uint64 `json:"id"`
	Username       string `gorm:"index:idx_username,unique;type:varchar(63);not null" json:"username"`
	Password       string `gorm:"type:varchar(255);not null" json:"password"`
	FollowingCount uint64 `gorm:"default:0;not null" json:"following_count"`
	FollowerCount  uint64 `gorm:"default:0;not null" json:"follower_count"`
	Avatar         string `gorm:"type:varchar(255);not null" json:"avatar"`
	WorkCount      uint64 `gorm:"default:0;not null" json:"work_count"`
	FavoriteCount  uint64 `gorm:"default:0;not null" json:"favorite_count"`
}

func (n *User) TableName() string {
	return constant.UserTableName
}

func CreateUser(user *User) (userID int64, err error) {
	if err := DB.Create(user).Error; err != nil {
		return 0, err
	}
	return int64(user.ID), nil
}

func SelectUserByID(userID uint64) (*User, error) {
	res := User{ID: userID}
	if err := DB.First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func SelectUserByName(username string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.Where("username = ?", username).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func IncreaseUserFavoriteCount(userID uint64) (uint64, error) {
	user := &User{
		ID: userID,
	}
	err := DB.First(&user).Error
	if err != nil {
		return 0, err
	}
	if err := DB.Model(&user).Update("favorite_count", user.FavoriteCount+1).Error; err != nil {
		return 0, err
	}
	return user.FavoriteCount, nil
}

func DecreaseUserFavoriteCount(userID uint64) (uint64, error) {
	user := &User{
		ID: userID,
	}
	err := DB.First(&user).Error
	if err != nil {
		return 0, err
	}
	if err := DB.Model(&user).Update("favorite_count", user.FavoriteCount-1).Error; err != nil {
		return 0, err
	}
	return user.FavoriteCount, nil
}
