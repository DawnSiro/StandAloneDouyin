package db

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	UserId   uint64 `json:"user_id"`
	ToUserId uint64 `json:"to_user_id"`
}

func (n *Relation) TableName() string {
	return constant.RelationTableName
}

func IsFollow(userId uint64, toUserId uint64) bool {
	relation := &Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}

	//follow by myself put false
	if userId == toUserId {
		return true
	}

	//other
	result := DB.Find(&relation, "user_id = ? and to_user_id = ?", userId, toUserId)

	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func AddFollow(userId uint64, toUserId uint64) error {
	relation := &Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}
	err := DB.Create(relation).Error
	return err
}

func DelFollow(userId uint64, toUserId uint64) error {
	relation := &Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}

	err := DB.Unscoped().Where("user_id = ? and to_user_id = ?", userId, toUserId).Delete(relation).Error
	return err
}

func GetFollowList(userID uint64) ([]*api.User, error) {
	commonResult := new([]*Relation)

	result := DB.Where("user_id = ?", userID).Find(&commonResult)

	results := make([]*api.User, 0)

	if result.RowsAffected == 0 {
		return results, nil
	}
	for i := 0; i < len(*commonResult); i++ {
		con1, err := SelectUserByUserID(uint((*commonResult)[i].ToUserId))
		if err != nil {
			return nil, err
		}

		results = append(results,
			&api.User{
				ID:            con1.ID,
				Name:          con1.Name,
				FollowCount:   con1.FollowCount,
				FollowerCount: con1.FollowerCount,
				IsFollow:      IsFollow((*commonResult)[i].ToUserId, userID),
				Avatar:        con1.Avatar,
			})
	}
	return results, nil
}

func GetFollowerList(userID uint64) ([]*api.User, error) {
	commonResult := new([]*Relation)

	result := DB.Where("to_user_id = ?", userID).Find(&commonResult)

	results := make([]*api.User, 0)

	if result.RowsAffected == 0 {
		return results, nil
	}
	for i := 0; i < len(*commonResult); i++ {
		con1, err := SelectUserByUserID(uint((*commonResult)[i].ToUserId))
		if err != nil {
			return nil, err
		}

		results = append(results,
			&api.User{
				ID:            con1.ID,
				Name:          con1.Name,
				FollowCount:   con1.FollowCount,
				FollowerCount: con1.FollowerCount,
				IsFollow:      IsFollow((*commonResult)[i].ToUserId, userID),
				Avatar:        con1.Avatar,
			})
	}
	return results, nil
}

func GetFriendList(userID uint64) ([]*api.FriendUser, error) {
	commonResult := new([]*Relation)

	result := DB.Where("user_id = ?", userID).Find(&commonResult)

	results := make([]*api.FriendUser, 0)
	messageResult := &FriendMessageResp{}

	if result.RowsAffected == 0 {
		return results, nil
	}
	for i := 0; i < len(*commonResult); i++ {
		//查看对方是否是自己的粉丝
		commonResult2 := new([]*Relation)
		con0 := DB.Where("user_id = ? and to_user_id = ?", (*commonResult)[i].ToUserId, userID).First(&commonResult2)
		if con0.RowsAffected == 0 {
			continue
		}

		con1, err := SelectUserByUserID(uint((*commonResult)[i].ToUserId))
		if err != nil {
			return nil, err
		}

		messageResult, err = GetLatestMsg(userID, (*commonResult)[i].ToUserId)
		if err != nil {
			return nil, err
		}

		results = append(results,
			&api.FriendUser{
				ID:            con1.ID,
				Name:          con1.Name,
				FollowCount:   con1.FollowCount,
				FollowerCount: con1.FollowerCount,
				IsFollow:      IsFollow((*commonResult)[i].ToUserId, userID),
				Avatar:        con1.Avatar,
				Message:       &messageResult.Content,
				MsgType:       int64(messageResult.MsgType),
			})
	}
	return results, nil
}
