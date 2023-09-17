package global

import (
	"sync"

	"douyin/pkg/config"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	ConfigPath                 string             // 配置文件路径
	Config                     config.System      // 系统配置信息
	FileTypeMap                sync.Map           // 文件前缀到文件类型的 map，使用 sync.Map 来保证并发安全
	DB                         *gorm.DB           // 数据库接口
	VideoRC                    *redis.Client      // 视频相关信息
	CommentRC                  *redis.Client      // 视频评论列表
	UserRC                     *redis.Client      // 用户信息
	MessageRC                  *redis.Client      // 消息聊天相关信息
	UserIDBloomFilter          *bloom.BloomFilter // 存储用户ID的布隆过滤器
	VideoIDBloomFilter         *bloom.BloomFilter // 存储视频ID的布隆过滤器
	CommentLuaScriptHash       string
	DeleteCommentLuaScriptHash string
	FollowLuaScriptHash        string
	CancelFollowLuaScriptHash  string
	UnLockLuaScriptHash        string
	PulsarClient               pulsar.Client // Pulsar 消息队列客户端
	FileSuffixWhiteList        = map[string]bool{".mp4": true}
)
