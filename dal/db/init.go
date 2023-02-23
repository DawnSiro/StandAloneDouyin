package db

import (
	"log"
	"os"
	"time"

	"douyin/pkg/constant"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var VideoCRC *redis.Client   // 视频评论列表
var VideoFRC *redis.Client   // 视频点赞
var UserInfoRC *redis.Client // 用户信息

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      logger.Info,            // 日志级别
			Colorful:      false,                  // 禁用彩色打印
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(constant.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger:                 newLogger,
		},
	)
	if err != nil {
		panic(err)
	}

	VideoCRC = redis.NewClient(&redis.Options{
		Addr: constant.RedisAddress,
		DB:   constant.VideoCRDB,
	})
	_, err = VideoCRC.Ping().Result()
	if err != nil {
		panic(err)
	}

	VideoFRC = redis.NewClient(&redis.Options{
		Addr: constant.RedisAddress,
		DB:   constant.VideoFRDB,
	})

	_, err = VideoFRC.Ping().Result()
	if err != nil {
		panic(err)
	}

	UserInfoRC = redis.NewClient(&redis.Options{
		Addr: constant.RedisAddress,
		DB:   constant.UserInfoRDB,
	})

	_, err = UserInfoRC.Ping().Result()
	if err != nil {
		panic(err)
	}

}
