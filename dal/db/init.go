package db

import (
	"douyin/constant"

	"github.com/go-redis/redis"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(constant.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		//Password: "123456",
		DB: 0,
	})
	_, err = RDB.Ping().Result()
	if err != nil {
		panic(err)
	}

}
