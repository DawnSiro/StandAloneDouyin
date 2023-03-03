package db

import (
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
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

// IsFollow ID 为 userID 的用户是否关注了 ID 为 toUserID 的用户
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
	result := global.DB.Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
		userID, toUserID, constant.DataNotDeleted).Limit(1).Find(&Relation{})
	if result.RowsAffected == 1 {
		return true
	}
	// 查询出错和没有数据都返回 false
	return false
}

// IsFriend ID 为 userID 的用户是否是 ID 为 toUserID 的用户的好友（互相关注）
func IsFriend(userID uint64, toUserID uint64) bool {
	// 默认不能给自己发消息
	if userID == 0 || toUserID == 0 || userID == toUserID {
		return false
	}

	// limit 2 需要使用切片而非单个结构体
	relation := make([]*Relation, 0, 2)
	result := global.DB.Where("(user_id = ? AND to_user_id = ? OR user_id = ? AND to_user_id = ?) AND is_deleted = ?",
		userID, toUserID, toUserID, userID, constant.DataNotDeleted).Limit(2).Find(&relation)
	if result.RowsAffected == 2 {
		return true
	}
	// 查询出错和没有数据都返回 false
	return false
}

func Follow(userID uint64, toUserID uint64) error {
	if userID == 0 || toUserID == 0 {
		hlog.Error("db.relation.Follow err:", errno.UserRequestParameterError)
		return errno.UserRequestParameterError
	}
	relation := &Relation{
		UserID:   userID,
		ToUserID: toUserID,
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 新增自己的关注数
		self := &User{ID: userID}
		err := tx.Select("following_count").First(self).Error
		if err != nil {
			return err
		}
		err = tx.Model(self).Update("following_count", self.FollowingCount+1).Error
		if err != nil {
			return err
		}
		// 新增关注用户的粉丝数
		opposite := &User{ID: toUserID}
		err = tx.Select("follower_count").First(opposite).Error
		if err != nil {
			return err
		}
		err = tx.Model(opposite).Update("follower_count", opposite.FollowerCount+1).Error
		if err != nil {
			return err
		}
		// 更新关注的关系
		// 先查询是否存在软删除的关注数据
		result := tx.Model(&Relation{}).Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
			userID, toUserID, constant.DataDeleted).Limit(1).Find(relation)
		// 如果有则修改为未删除
		if result.RowsAffected == 1 {
			return tx.Model(relation).Update("is_deleted", constant.DataNotDeleted).Error
		}
		// 没有则新建
		return tx.Create(relation).Error
	})
}

func CancelFollow(userID uint64, toUserID uint64) error {
	if userID == 0 || toUserID == 0 {
		return errno.UserRequestParameterError
	}

	relation := &Relation{}
	result := global.DB.Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
		userID, toUserID, constant.DataNotDeleted).Limit(1).Find(relation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.UserRequestParameterError
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 减少自己的关注数
		self := &User{ID: userID}
		err := tx.Select("following_count").First(self).Error
		if err != nil {
			return err
		}
		err = tx.Model(self).Update("following_count", self.FollowingCount-1).Error
		if err != nil {
			return err
		}
		// 减少关注用户的粉丝数
		opposite := &User{ID: toUserID}
		err = tx.Select("follower_count").First(opposite).Error
		if err != nil {
			return err
		}
		err = tx.Model(opposite).Update("follower_count", opposite.FollowerCount-1).Error
		if err != nil {
			return err
		}

		// 去除关注的关系
		return global.DB.Model(relation).Where("is_deleted = ?", constant.DataNotDeleted).
			Update("is_deleted", constant.DataDeleted).Error
	})

}

func GetFollowList(userID uint64) ([]*User, error) {
	relations := make([]*Relation, 0)
	res := make([]*User, 0)

	result := global.DB.Where("user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&relations)
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

	result := global.DB.Where("to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&relations)
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

func GetFriendList(userID uint64) ([]*User, error) {
	friendUserID := make([]int64, 0)
	res := make([]*User, 0)

	// 查询关注自己的所有粉丝的 userID
	followerQuery := global.DB.Select("user_id").Table(constant.RelationTableName).
		Where("to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted)
	// 查询自己关注的，并且在自己的粉丝 userID 集合里的用户
	err := global.DB.Select("to_user_id").Table(constant.RelationTableName).
		Where("user_id = ? AND is_deleted = ? AND to_user_id IN (?)",
			userID, constant.DataNotDeleted, followerQuery).Find(&friendUserID).Error
	if err != nil {
		return nil, err
	}
	// 没有好友，直接返回
	if len(friendUserID) == 0 {
		return res, nil
	}
	err = global.DB.Where("id IN (?)", friendUserID).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

type RelationUserData struct {
	UID            uint64 `gorm:"column:uid"`
	Username       string
	FollowingCount uint64
	FollowerCount  uint64
	Avatar         string
	IsFollow       bool
}

func SelectFollowDataListByUserID(userID uint64) ([]*RelationUserData, error) {
	res := make([]*RelationUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow", constant.DataNotDeleted).Table("`user` AS u").
		Joins("RIGHT JOIN relation AS r ON r.to_user_id = u.id").Where("r.user_id = ?", userID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectFollowerDataListByUserID(userID uint64) ([]*RelationUserData, error) {
	res := make([]*RelationUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow", constant.DataNotDeleted).Table("`user` AS u").
		Joins("RIGHT JOIN relation AS r ON r.user_id = u.id").Where("r.to_user_id = ?", userID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
