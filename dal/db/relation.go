package db

import (
	"douyin/biz/model/api"
	"douyin/pkg/constant"
	"errors"
)

type Relation struct {
	ID        uint64 `json:"id"`
	IsDeleted uint8  `gorm:"default:0;not null" json:"is_deleted"`
	UserID    uint64 `gorm:"not null" json:"user_id"`
	ToUserID  uint64 `gorm:"not null" json:"to_user_id"`
}

func (n *Relation) TableName() string {
	return constant.RelationTableName
}

func IsFollow(userID uint64, toUserID uint64) bool {
	// 未登录默认未关注
	if userID == 0 || toUserID == 0 {
		return false
	}

	// 默认自己关注自己，看自己的视频时头像上没有关注的那个加号按钮
	if userID == toUserID {
		return true
	}

	// 查不到关注记录则为未关注
	result := DB.Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
		userID, toUserID, constant.DataNotDeleted).Limit(1).Find(&Relation{})
	if result.RowsAffected == 1 {
		return true
	}
	// 查询出错和没有数据都返回 false
	return false
}

func Follow(userID uint64, toUserID uint64) error {
	if userID == 0 || toUserID == 0 {
		return errors.New("follow failed")
	}
	relation := &Relation{
		UserID:   userID,
		ToUserID: toUserID,
	}
	// 先查询是否存在软删除的关注数据
	result := DB.Where("user_id = ? AND to_user_id = ? AND AND is_deleted = ?",
		userID, toUserID, constant.DataNotDeleted).Limit(1).Find(relation)
	// 如果有则修改为未删除
	if result.RowsAffected == 1 {
		return DB.Model(relation).Update("is_deleted", constant.DataDeleted).Error
	}
	//查询该用户是否存在
	result = DB.Table("user").Where("id = ?", toUserID).Find(&api.User{})
	if result.RowsAffected == 0 {
		return errors.New("关注用户不存在")
	}
	// 没有则新建
	return DB.Create(relation).Error
}

func CancelFollow(userID uint64, toUserID uint64) error {
	if userID == 0 || toUserID == 0 {
		return errors.New("delete data failed")
	}

	relation := &Relation{
		UserID:   userID,
		ToUserID: toUserID,
	}

	err := DB.Model(relation).Where("user_id = ? AND to_user_id = ?",
		userID, toUserID).Update("is_deleted", constant.DataDeleted)
	if err.RowsAffected == 0 {
		return errors.New("请勿重复操作")
	}
	return nil
}

func GetFollowList(userID uint64) ([]*User, error) {
	relations := make([]*Relation, 0)
	res := make([]*User, 0)

	result := DB.Where("user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&relations)
	if result.RowsAffected == 0 {
		return res, nil
	}

	for i := 0; i < len(relations); i++ {
		u, err := SelectUserByID(relations[i].ToUserID)
		if err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	return res, nil
}

func GetFollowerList(userID uint64) ([]*User, error) {
	relations := make([]*Relation, 0)
	res := make([]*User, 0)

	result := DB.Where("to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&relations)
	if result.RowsAffected == 0 {
		return res, nil
	}

	for i := 0; i < len(relations); i++ {
		u, err := SelectUserByID(relations[i].UserID)
		if err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	return res, nil
}

func GetFriendList(userID uint64) ([]*api.FriendUser, error) {
	relations := make([]*Relation, 0)
	results := make([]*api.FriendUser, 0)

	result := DB.Where("user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&relations)
	if result.RowsAffected == 0 {
		return results, nil
	}

	for i := 0; i < len(relations); i++ {
		//查看对方是否是自己的粉丝
		rs2 := make([]*Relation, 0)
		result := DB.Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
			relations[i].ToUserID, userID, constant.DataNotDeleted).Limit(1).Find(&rs2)
		if result.RowsAffected == 0 {
			continue
		}

		u, err := SelectUserByID(relations[i].ToUserID)
		if err != nil {
			return nil, err
		}

		messageResult, err := GetLatestMsg(userID, relations[i].ToUserID)
		if err != nil {
			return nil, err
		}

		followCount := int64(u.FollowingCount)
		followerCount := int64(u.FollowerCount)
		results = append(results,
			&api.FriendUser{
				ID:            int64(u.ID),
				Name:          u.Username,
				FollowCount:   &followCount,
				FollowerCount: &followerCount,
				IsFollow:      IsFollow(relations[i].ToUserID, userID),
				Avatar:        u.Avatar,
				Message:       &messageResult.Content,
				MsgType:       int8(messageResult.MsgType),
			})
	}

	return results, nil
}
