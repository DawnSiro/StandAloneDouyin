package global

import (
	"douyin/pkg/config"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Config     config.System // 系统配置信息
	DB         *gorm.DB      // 数据库接口
	VideoCRC   *redis.Client // 视频评论列表
	VideoFRC   *redis.Client // 视频点赞
	UserInfoRC *redis.Client // 用户信息
)
