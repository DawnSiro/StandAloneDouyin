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

var VideoCRDB *redis.Client // 视频评论列表
var VideoFRDB *redis.Client // 视频点赞

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

	// TODO 使用 Viper 从配置文件中读取
	VideoCRDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	_, err = VideoCRDB.Ping().Result()
	if err != nil {
		panic(err)
	}

	VideoFRDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	})

	_, err = VideoFRDB.Ping().Result()
	if err != nil {
		panic(err)
	}

}
