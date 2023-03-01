package initialize

import (
	"fmt"

	"douyin/pkg/global"

	"github.com/go-redis/redis"
)

func Redis() {
	// VideoCRC链接
	global.VideoCRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.VideoCRCRedisConfig.Host, global.Config.VideoCRCRedisConfig.Port),
		Password: global.Config.VideoCRCRedisConfig.Password,
		DB:       global.Config.VideoCRCRedisConfig.DB,
		PoolSize: global.Config.VideoCRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.VideoCRC.Ping().Result(); err != nil {
		panic(err.Error())
	}

	// VideoFRC链接
	global.VideoFRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.VideoFRCRedisConfig.Host, global.Config.VideoFRCRedisConfig.Port),
		Password: global.Config.VideoFRCRedisConfig.Password,
		DB:       global.Config.VideoFRCRedisConfig.DB,
		PoolSize: global.Config.VideoFRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.VideoFRC.Ping().Result(); err != nil {
		panic(err.Error())
	}

	// UserInfoRC链接
	global.UserInfoRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.UserInfoRCRedisConfig.Host, global.Config.UserInfoRCRedisConfig.Port),
		Password: global.Config.UserInfoRCRedisConfig.Password,
		DB:       global.Config.UserInfoRCRedisConfig.DB,
		PoolSize: global.Config.UserInfoRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.UserInfoRC.Ping().Result(); err != nil {
		panic(err.Error())
	}
}
