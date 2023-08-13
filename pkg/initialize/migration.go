package initialize

import "douyin/pkg/global"
import "douyin/dal/db"

func migration() {
	//自动迁移模式
	global.DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&db.Comment{}, &db.CommentData{}, &db.UserFavoriteVideo{}, &db.Message{},
			&db.Relation{}, &db.RelationUserData{}, &db.User{},
			&db.Video{}, &db.VideoData{})

}
