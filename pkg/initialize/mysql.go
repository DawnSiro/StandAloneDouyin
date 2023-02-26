package initialize

import (
	"douyin/pkg/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func MySQL() {

	username := global.Config.MySQLConfig.Username // 账号
	password := global.Config.MySQLConfig.Password // 密码
	host := global.Config.MySQLConfig.Host         // 数据库地址，可以是Ip或者域名
	port := global.Config.MySQLConfig.Port         // 数据库端口
	dbName := global.Config.MySQLConfig.DBname     // 数据库名
	// dsn := "用户名:密码@tcp(地址:端口)/数据库名"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)

	// 配置Gorm连接到MySQL
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 0, // 慢 SQL 阈值
			LogLevel:                  logger.Info,          // 日志级别
			IgnoreRecordNotFoundError: true,                 // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                // 禁用彩色打印
		},
	)
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: newLogger,
	}); err == nil {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(global.Config.MySQLConfig.MaxOpenConns) // 设置数据库最大连接数
		sqlDB.SetMaxIdleConns(global.Config.MySQLConfig.MaxIdleConns) // 设置上数据库最大闲置连接数
		global.DB = db
	} else {
		panic("connect server failed")
	}

}
