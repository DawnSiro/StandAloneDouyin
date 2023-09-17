package model

import "douyin/pkg/constant"

type User struct {
	ID              uint64 `json:"id"`                                                                  // 自增主键
	Username        string `gorm:"index:idx_username,unique;type:varchar(63);not null" json:"username"` // 用户名
	Password        string `gorm:"type:varchar(255);not null" json:"password"`                          // 用户密码
	FollowingCount  int64  `gorm:"default:0;not null" json:"following_count"`                           // 关注数
	FollowerCount   int64  `gorm:"default:0;not null" json:"follower_count"`                            // 粉丝数
	Avatar          string `gorm:"type:varchar(255);null" json:"avatar"`                                // 用户头像
	BackgroundImage string `gorm:"type:varchar(255);not null" json:"background_image"`                  // 用户个人页顶部大图
	Signature       string `gorm:"type:varchar(255);not null" json:"signature"`                         // 个人简介
	TotalFavorited  int64  `gorm:"default:0;not null" json:"total_favorited"`                           // 获赞数量
	WorkCount       int64  `gorm:"default:0;not null" json:"work_count"`                                // 作品数量
	FavoriteCount   int64  `gorm:"default:0;not null" json:"favorite_count"`                            // 点赞数量
}

func (u *User) TableName() string {
	return constant.UserTableName
}
