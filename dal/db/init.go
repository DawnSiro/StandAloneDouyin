package db

import (
	"douyin/pkg/constant"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

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

	RDB = redis.NewClient(&redis.Options{
		//Addr: "172.17.0.1:6379",
		Addr: "127.0.0.1:6379",
		//Password: "123456",
		DB: 0,
	})
	_, err = RDB.Ping().Result()
	if err != nil {
		panic(err)
	}

}
