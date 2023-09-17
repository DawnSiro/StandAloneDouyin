package db

import (
	"douyin/dal/model"
	"douyin/pkg/global"
)

func CreateUser(user *model.User) (userID uint64, err error) {
	// 这里不指明更新的字段的话，会全部赋零值，就无法使用数据库的默认值了
	if err := global.DB.Select("username", "password").Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func SelectUserByID(userID uint64) (*model.User, error) {
	// TODO 优化查询字段
	res := model.User{ID: userID}
	if err := global.DB.First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func SelectUserByIDs(ids map[uint64]struct{}) ([]*model.User, error) {
	res := make([]*model.User, 0)
	err := global.DB.Where("id IN (?)", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectUserByName(username string) (*model.User, error) {
	res := &model.User{}
	if err := global.DB.Where("username = ?", username).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func IncreaseUserFavoriteCount(userID uint64) (int64, error) {
	user := &model.User{ID: userID}
	err := global.DB.First(&user).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&user).Update("favorite_count", user.FavoriteCount+1).Error; err != nil {
		return 0, err
	}
	return user.FavoriteCount, nil
}

func DecreaseUserFavoriteCount(userID uint64) (int64, error) {
	user := &model.User{ID: userID}
	err := global.DB.First(&user).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&user).Update("favorite_count", user.FavoriteCount-1).Error; err != nil {
		return 0, err
	}
	return user.FavoriteCount, nil
}

func SelectTopFollowers(limit int) ([]*model.User, error) {
	res := make([]*model.User, 0)
	err := global.DB.Order("follower_count DESC").Limit(limit).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
