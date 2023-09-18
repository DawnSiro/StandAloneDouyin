package db

import (
	"time"

	"douyin/dal/model"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

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
		userID, toUserID, constant.DataNotDeleted).Limit(1).Find(&model.Relation{})
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
	relation := make([]*model.Relation, 0, 2)
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
	relation := &model.Relation{
		UserID:      userID,
		ToUserID:    toUserID,
		CreatedTime: time.Now(),
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 新增自己的关注数
		self := &model.User{ID: userID}
		err := tx.Select("username, avatar, following_count").First(self).Error
		if err != nil {
			return err
		}
		err = tx.Model(self).Update("following_count", self.FollowingCount+1).Error
		if err != nil {
			return err
		}
		// 新增关注用户的粉丝数
		opposite := &model.User{ID: toUserID}
		err = tx.Select("username, avatar, follower_count").First(opposite).Error
		if err != nil {
			return err
		}
		err = tx.Model(opposite).Update("follower_count", opposite.FollowerCount+1).Error
		if err != nil {
			return err
		}
		// 更新关注的关系
		// 先查询是否存在软删除的关注数据
		result := tx.Model(&model.Relation{}).Where("user_id = ? AND to_user_id = ? AND is_deleted = ?",
			userID, toUserID, constant.DataDeleted).Limit(1).Find(relation)
		// 如果有则修改为未删除
		if result.RowsAffected == 1 {
			return tx.Model(relation).Update("is_deleted", constant.DataNotDeleted).Error
		}
		// 没有则新建
		err = tx.Create(relation).Error
		if err != nil {
			return err
		}
		// 更新缓存
		return rdb.Follow(&model.FanUserRedisData{
			UID:      self.ID,
			Username: self.Username,
			Avatar:   self.Avatar,
		}, &model.FollowUserData{
			UID:         opposite.ID,
			Username:    opposite.Username,
			Avatar:      opposite.Avatar,
			CreatedTime: relation.CreatedTime,
		})
	})
}

func CancelFollow(userID uint64, toUserID uint64) error {
	if userID == 0 || toUserID == 0 {
		return errno.UserRequestParameterError
	}

	relation := &model.Relation{}
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
		self := &model.User{ID: userID}
		err := tx.Select("following_count").First(self).Error
		if err != nil {
			return err
		}
		err = tx.Model(self).Update("following_count", self.FollowingCount-1).Error
		if err != nil {
			return err
		}
		// 减少关注用户的粉丝数
		opposite := &model.User{ID: toUserID}
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

// GetFollowList 获取关注列表用户信息
func GetFollowList(userID uint64) ([]*model.User, error) {
	relations := make([]*model.Relation, 0)
	res := make([]*model.User, 0)

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

func GetFollowerList(userID uint64) ([]*model.User, error) {
	relations := make([]*model.Relation, 0)
	res := make([]*model.User, 0)

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

func GetFriendList(userID uint64) ([]*model.User, error) {
	friendUserID := make([]int64, 0)
	res := make([]*model.User, 0)

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

// SelectFollowUserListByUserID 查询关注列表的用户信息
// 这里不用 IN 是因为 IN 的性能没有 JOIN 好
func SelectFollowUserListByUserID(userID uint64) ([]*model.FollowUserData, error) {
	res := make([]*model.FollowUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.avatar, r.created_time").
		Table("`user` AS u").
		Joins("JOIN relation AS r ON r.to_user_id = u.id").
		Where("r.user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SelectFanUserListByUserID 查询粉丝列表的用户信息
// 这里不能用 IN 是因为 IN 的性能没有 JOIN 好
func SelectFanUserListByUserID(userID uint64) ([]*model.FanUserData, error) {
	res := make([]*model.FanUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.avatar, r.created_time").
		Table("`user` AS u").
		Joins("JOIN relation AS r ON r.user_id = u.id").
		Where("r.to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectFollowDataListByUserID(userID uint64) ([]*model.RelationUserData, error) {
	res := make([]*model.RelationUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow", constant.DataNotDeleted).Table("`user` AS u").
		Joins("RIGHT JOIN relation AS r ON r.to_user_id = u.id").Where("r.user_id = ?", userID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectFollowerDataListByUserID(userID uint64) ([]*model.RelationUserData, error) {
	res := make([]*model.RelationUserData, 0)
	err := global.DB.Select("u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow", constant.DataNotDeleted).Table("`user` AS u").
		Joins("RIGHT JOIN relation AS r ON r.user_id = u.id").Where("r.to_user_id = ?", userID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SelectFollowUserZSet 查询一个用户所有的关注列表，用于存储到 redis 中
func SelectFollowUserZSet(userID uint64) ([]*rdb.RelationZSet, error) {
	followUserList := make([]*rdb.RelationZSet, 0)
	err := global.DB.Select("id, to_user_id AS member").Table(constant.RelationTableName).
		Where("user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&followUserList).Error
	if err != nil {
		return nil, err
	}
	return followUserList, nil
}

// SelectFanUserZSet 查询一个用户所有的粉丝列表，用于存储到 redis 中
func SelectFanUserZSet(userID uint64) ([]*rdb.RelationZSet, error) {
	fanUserList := make([]*rdb.RelationZSet, 0)
	err := global.DB.Select("id, user_id AS member").Table(constant.RelationTableName).
		Where("to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&fanUserList).Error
	if err != nil {
		return nil, err
	}
	return fanUserList, nil
}

// SelectFollowUserIDSet 查询用户关注的 用户的ID 集合
func SelectFollowUserIDSet(userID uint64) ([]uint64, error) {
	followUserList := make([]uint64, 0)
	err := global.DB.Select("to_user_id").Table(constant.RelationTableName).
		Where("user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&followUserList).Error
	if err != nil {
		return nil, err
	}
	return followUserList, nil
}

func SelectFanUserIDSet(userID uint64) ([]uint64, error) {
	fanUserList := make([]uint64, 0)
	err := global.DB.Select("user_id").Table(constant.RelationTableName).
		Where("to_user_id = ? AND is_deleted = ?", userID, constant.DataNotDeleted).Find(&fanUserList).Error
	if err != nil {
		return nil, err
	}
	return fanUserList, nil
}

// SelectFriendDataList 获取好友列表
func SelectFriendDataList(userID uint64) ([]*model.FriendUserData, error) {
	friendUserID := make([]int64, 0)
	uList := make([]*model.User, 0)

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
		return nil, nil
	}
	err = global.DB.Where("id IN (?)", friendUserID).Find(&uList).Error
	if err != nil {
		return nil, err
	}
	uMap := make(map[uint64]*model.User)
	for i := 0; i < len(uList); i++ {
		uMap[uList[i].ID] = uList[i]
	}

	// 查询最新消息
	userSendMessage := make([]*model.Message, len(friendUserID))
	err = global.DB.Raw("SELECT m.* "+
		"FROM message AS m "+
		"JOIN (SELECT from_user_id, MAX(created_time) AS max_created_time "+
		"FROM message "+
		"WHERE to_user_id IN (?) AND from_user_id = ? "+
		"GROUP BY from_user_id) AS sub "+
		"ON m.from_user_id = sub.from_user_id AND m.created_time = sub.max_created_time "+
		"ORDER BY m.created_time DESC;",
		friendUserID, userID).Scan(&userSendMessage).Error
	if err != nil {
		return nil, err
	}

	userReceiveMessage := make([]*model.Message, len(friendUserID))
	err = global.DB.Raw("SELECT m.* "+
		"FROM message AS m "+
		"JOIN (SELECT from_user_id, MAX(created_time) AS max_created_time "+
		"FROM message "+
		"WHERE to_user_id = ? AND from_user_id IN (?) "+
		"GROUP BY from_user_id) AS sub "+
		"ON m.from_user_id = sub.from_user_id AND m.created_time = sub.max_created_time "+
		"ORDER BY m.created_time DESC;",
		userID, friendUserID).Scan(&userReceiveMessage).Error
	if err != nil {
		return nil, err
	}

	// 拼出一个结果然后返回
	set := make(map[uint64]struct{}, len(friendUserID))
	res := make([]*model.FriendUserData, 0, len(friendUserID))
	// 合并两个集合
	for _, message := range userSendMessage {
		set[message.ToUserID] = struct{}{}
		res = append(res, &model.FriendUserData{
			ID:       message.ToUserID,
			Name:     uMap[message.ToUserID].Username,
			IsFollow: true, // 都是朋友了当然是关注了
			Avatar:   uMap[message.ToUserID].Avatar,
			Message:  message.Content,
			MsgType:  1, // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
		})
		for _, m := range userReceiveMessage {
			if message.ToUserID == m.FromUserID {
				// 如果另一个集合里的时间更晚则更新
				if message.CreatedTime.Before(m.CreatedTime) {
					res[len(res)-1].Message = m.Content
					res[len(res)-1].MsgType = 0 // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
				}
			}
		}
	}

	for _, message := range userReceiveMessage {
		// 没有加入结果集的加入一下
		if _, ok := set[message.FromUserID]; !ok {
			res = append(res, &model.FriendUserData{
				ID:       message.FromUserID,
				Name:     uMap[message.FromUserID].Username,
				IsFollow: true, // 都是朋友了当然是关注了
				Avatar:   uMap[message.FromUserID].Avatar,
				Message:  message.Content,
				// message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
				MsgType: 0,
			})
		}
	}

	return res, nil
}
