package db

import (
	"douyin/constants"
	"errors"
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	UserId   uint64 `json:"user_id"`
	ToUserId uint64 `json:"to_user_id"`
}

func (n *Relation) TableName() string {
	return constants.RelationTableName
}

func IsFollow(userId uint64, toUserId uint64) bool {
	relation := &Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}

	//follow by myself put false
	if userId == toUserId {
		return false
	}

	//other
	result := DB.First(&relation, "user_id = ? and to_user_id = ?", userId, toUserId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
