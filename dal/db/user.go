package db

import (
	"douyin/pkg/constant"
	"douyin/pkg/global"
)

type User struct {
	ID              uint64 `json:"id"`                                                                  // 自增主键
	Username        string `gorm:"index:idx_username,unique;type:varchar(63);not null" json:"username"` //用户名
	Password        string `gorm:"type:varchar(255);not null" json:"password"`                          //用户密码
	FollowingCount  int64  `gorm:"default:0;not null" json:"following_count"`                           //关注数
	FollowerCount   int64  `gorm:"default:0;not null" json:"follower_count"`                            //粉丝数
	Avatar          string `gorm:"type:varchar(255);not null" json:"avatar"`                            //用户头像
	BackgroundImage string `gorm:"type:varchar(255);not null" json:"background_image"`                  //用户个人页顶部大图
	Signature       string `gorm:"type:varchar(255);not null" json:"signature"`                         //个人简介
	TotalFavorited  int64  `gorm:"default:0;not null" json:"total_favorited"`                           //获赞数量
	WorkCount       int64  `gorm:"default:0;not null" json:"work_count"`                                //作品数量
	FavoriteCount   int64  `gorm:"default:0;not null" json:"favorite_count"`                            //点赞数量
}

func (u *User) TableName() string {
	return constant.UserTableName
}

func CreateUser(user *User) (userID uint64, err error) {
	// 这里不指明更新的字段的话，会全部赋零值，就无法使用数据库的默认值了
	if err := global.DB.Select("username", "password").Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func SelectUserByID(userID uint64) (*User, error) {
	res := User{ID: userID}
	if err := global.DB.First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func SelectUserByIDs(ids map[uint64]struct{}) ([]*User, error) {
	res := make([]*User, 0)
	err := global.DB.Where("id IN (?)", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectUserByName(username string) ([]*User, error) {
	res := make([]*User, 0)
	if err := global.DB.Where("username = ?", username).Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func IncreaseUserFavoriteCount(userID uint64) (int64, error) {
	user := &User{ID: userID}
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
	user := &User{ID: userID}
	err := global.DB.First(&user).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&user).Update("favorite_count", user.FavoriteCount-1).Error; err != nil {
		return 0, err
	}
	return user.FavoriteCount, nil
}
