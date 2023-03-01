package global

import (
	"sync"

	"douyin/pkg/config"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Config              config.System // 系统配置信息
	FileTypeMap         sync.Map      // 文件前缀到文件类型的 map，使用 sync.Map 来保证并发安全
	DB                  *gorm.DB      // 数据库接口
	VideoCRC            *redis.Client // 视频评论列表
	VideoFRC            *redis.Client // 视频点赞
	UserInfoRC          *redis.Client // 用户信息
	FileSuffixWhiteList = map[string]bool{".mp4": true}
)
