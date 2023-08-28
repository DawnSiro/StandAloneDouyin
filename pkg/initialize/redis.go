package initialize

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strings"
	"sync"

	"douyin/pkg/global"

	"github.com/go-redis/redis"
)

var cancelSubscription context.CancelFunc

func Redis() {
	ctx := context.Background()

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

	// Subscribe redis cache topics
	_, cancelSubscription = context.WithCancel(ctx)
	go SubscribeToTopicChanges()

}

func SubscribeToTopicChanges() {
	// Cancel the context when the Goroutine exits
	defer cancelSubscription()

	// Wait for all goroutines to finish
	var wg sync.WaitGroup

	topics := []string{"friendList_changes", "commentList_changes"}

	for _, topic := range topics {
		pubSub := global.UserInfoRC.Subscribe(topic)
		wg.Add(1)

		go func(topic string, pubSub *redis.PubSub) {
			defer wg.Done()

			defer func(pubSub *redis.PubSub) {
				err := pubSub.Close()
				if err != nil {
					hlog.Error("pkg.initialize.redis.SubscribeToChanges err: " + err.Error())
				}
			}(pubSub)

			// Handle message
			for {
				msg, err := pubSub.ReceiveMessage()
				if err != nil {
					hlog.Error("pkg.initialize.redis.SubscribeToChanges err: " + err.Error())
				}
				parts := strings.Split(msg.Payload, "&")

				if parts[0] == "friend_followed" || parts[0] == "friend_unfollowed" {
					// Invalidate the cached friend list
					friendListCacheKeyFromUser := "friendList:" + parts[1]
					friendListCacheKeyToUser := "friendList:" + parts[2]
					followerListCacheKeyFromUser := "followerList:" + parts[1]
					followerListCacheKeyToUser := "followerList:" + parts[2]
					followListCacheKeyFromUser := "followList:" + parts[1]
					followListCacheKeyToUser := "followList:" + parts[2]
					err := global.UserInfoRC.Del(friendListCacheKeyFromUser).Err()
					err = global.UserInfoRC.Del(friendListCacheKeyToUser).Err()
					err = global.UserInfoRC.Del(followerListCacheKeyFromUser).Err()
					err = global.UserInfoRC.Del(followerListCacheKeyToUser).Err()
					err = global.UserInfoRC.Del(followListCacheKeyFromUser).Err()
					err = global.UserInfoRC.Del(followListCacheKeyToUser).Err()
					if err != nil {
						hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: ", err.Error())
					}
				} else if parts[0] == "comment_added" || parts[0] == "comment_deleted" {
					// Invalidate the cache for comment
					videoComment := "commentList:" + parts[1] + ":" + parts[2]
					err := global.UserInfoRC.Del(videoComment).Err()
					if err != nil {
						hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.comment err: ", err.Error())
					}
				}
			}
		}(topic, pubSub)
	}

	wg.Wait()
}
