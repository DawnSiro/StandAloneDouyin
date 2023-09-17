package rdb

import (
	"douyin/pkg/viper"
	"flag"
	"fmt"
	"sync"

	"douyin/pkg/constant"
	"douyin/pkg/global"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
)

func InitRedis() {
	// 对需要用到的 Redis 客户端进行缓存
	initClient()

	// 初始化布隆过滤器
	initBloomFilter()

	// 对需要用到的 Lua 脚本进行缓存
	initLuaScript()

	// TODO 对可能访问量较高的数据（如头部Up)进行缓存预热
}

// initClient 对需要用到的 Redis 客户端进行缓存
func initClient() {
	// VideoRC链接
	global.VideoRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.VideoRCRedisConfig.Host, global.Config.VideoRCRedisConfig.Port),
		Password: global.Config.VideoRCRedisConfig.Password,
		DB:       global.Config.VideoRCRedisConfig.DB,
		PoolSize: global.Config.VideoRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.VideoRC.Ping().Result(); err != nil {
		panic(err.Error())
	}

	// VideoCRC链接
	global.CommentRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.CommentRCRedisConfig.Host, global.Config.CommentRCRedisConfig.Port),
		Password: global.Config.CommentRCRedisConfig.Password,
		DB:       global.Config.CommentRCRedisConfig.DB,
		PoolSize: global.Config.CommentRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.CommentRC.Ping().Result(); err != nil {
		panic(err.Error())
	}

	// UserRC 链接
	global.UserRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.UserRCRedisConfig.Host, global.Config.UserRCRedisConfig.Port),
		Password: global.Config.UserRCRedisConfig.Password,
		DB:       global.Config.UserRCRedisConfig.DB,
		PoolSize: global.Config.UserRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.UserRC.Ping().Result(); err != nil {
		panic(err.Error())
	}

	// Message 链接
	global.MessageRC = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.MessageRCRedisConfig.Host, global.Config.MessageRCRedisConfig.Port),
		Password: global.Config.MessageRCRedisConfig.Password,
		DB:       global.Config.MessageRCRedisConfig.DB,
		PoolSize: global.Config.MessageRCRedisConfig.PoolSize,
	})
	// 检查 Redis 连通性
	if _, err := global.MessageRC.Ping().Result(); err != nil {
		panic(err.Error())
	}
}

// initBloomFilter 初始化布隆过滤器
// TODO 考虑对布隆过滤器需要存储的数据进行预计算
func initBloomFilter() {
	global.UserIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	global.VideoIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
}

// initLuaScript 对需要用到的 Lua 脚本进行缓存
// 相同节点不同数据库共享缓存的 Lua 脚本，如果使用集群则需要手动或者使用工具进行同步
func initLuaScript() {
	var err error
	// 评论部分
	global.CommentLuaScriptHash, err = global.UserRC.ScriptLoad(constant.CommentLuaScript).Result()
	if err != nil {
		panic(err)
	}

	// 删除评论部分
	global.DeleteCommentLuaScriptHash, err = global.UserRC.ScriptLoad(constant.DeleteCommentLuaScript).Result()
	if err != nil {
		panic(err)
	}

	// 关注部分
	global.FollowLuaScriptHash, err = global.UserRC.ScriptLoad(constant.FollowLuaScript).Result()
	if err != nil {
		panic(err)
	}

	// 取消关注部分
	global.CancelFollowLuaScriptHash, err = global.UserRC.ScriptLoad(constant.CancelFollowLuaScript).Result()
	if err != nil {
		panic(err)
	}

	// 解锁部分，相同节点不同数据库共享 Lua 脚本
	global.UnLockLuaScriptHash, err = global.UserRC.ScriptLoad(constant.UnLockLuaScript).Result()
	if err != nil {
		panic(err)
	}
}

var once sync.Once

func InitTest() {
	once.Do(func() {
		flag.StringVar(&global.ConfigPath, "c", "../../pkg/config/config.yml", "config file path")
		flag.Parse()
		viper.InitConfig()
		InitRedis()
	})
}
