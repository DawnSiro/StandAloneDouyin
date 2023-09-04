package initialize

import (
	"context"
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"douyin/biz/service"
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

	// Start a goroutine to periodically warm up the cache
	go func() {
		for {
			err := WarmUpCacheForTopFollowers()
			if err != nil {
				hlog.Error("WarmUpCacheForTopFollowers error: " + err.Error())
			}

			err = WarmUpCacheForTopRelationData()
			if err != nil {
				hlog.Error("WarmUpCacheForTopRelationData error: " + err.Error())
			}

			err = WarmUpCacheForTopFavoriteVideos()
			if err != nil {
				hlog.Error("WarmUpCacheForTopFavoriteVideos error: " + err.Error())
			}

			err = WarmUpCacheForTopComments()
			if err != nil {
				hlog.Error("WarmUpCacheForTopComments error: " + err.Error())
			}

			// Wait for 15 minutes before the next warm-up
			time.Sleep(15 * time.Minute)
		}
	}()

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
				ID1, _ := strconv.ParseUint(parts[1], 10, 64)
				ID2, _ := strconv.ParseUint(parts[2], 10, 64)

				switch parts[0] {
				case "friend_followed":
					updateUserRelationships(ID1, ID2, true)
				case "friend_unfollowed":
					updateUserRelationships(ID1, ID2, false)
				case "comment_added", "comment_deleted":
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

func updateUserRelationships(followerID, followedID uint64, follow bool) {
	updateFollowerList(followerID, followedID, follow)

	updateFollowList(followerID, followedID, follow)

	updateFriendLists(followerID, followedID, follow)
}

func updateFollowerList(userID, followedID uint64, follow bool) {
	// Update the follower list for followedID

	followerListCacheKeyToUser := fmt.Sprintf("followList:%d", followedID)

	// Get the current follower list from Redis
	currentFollowerListJSON, err := global.UserInfoRC.Get(followerListCacheKeyToUser).Result()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: ", err.Error())
		return
	}

	var currentFollowerList []*db.RelationUserData
	if err := json.Unmarshal([]byte(currentFollowerListJSON), &currentFollowerList); err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error decoding follower list, ", err.Error())
		return
	}

	if follow {
		// If follow is true, add the user if not found in the list
		// TODO:这里需要构造结构体
		currentUserData := &db.RelationUserData{UID: userID}
		currentFollowerList = append(currentFollowerList, currentUserData)
	} else {
		// If follow is false, remove the user from the list if found
		updatedFollowerList := make([]*db.RelationUserData, 0, len(currentFollowerList)-1)
		for _, user := range currentFollowerList {
			if user.UID != userID {
				updatedFollowerList = append(updatedFollowerList, user)
			}
		}
		currentFollowerList = updatedFollowerList

	}

	// Update the follower list in Redis
	updatedFollowerListJSON, _ := json.Marshal(currentFollowerList)
	err = global.UserInfoRC.Set(followerListCacheKeyToUser, updatedFollowerListJSON, 0).Err()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error updating follower list, ", err.Error())
	}
}

func updateFollowList(userID, followedID uint64, follow bool) {
	// Update the follow list for userID

	followListCacheKey := fmt.Sprintf("followList:%d", userID)

	// Get the current follow list from Redis
	currentFollowListJSON, err := global.UserInfoRC.Get(followListCacheKey).Result()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: ", err.Error())
		return
	}

	var currentFollowList []*db.RelationUserData
	if err := json.Unmarshal([]byte(currentFollowListJSON), &currentFollowList); err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error decoding follow list, ", err.Error())
		return
	}

	if follow {
		// If follow is true, add the followed user if not found in the list
		// TODO:这里需要构造结构体
		followedUserData := &db.RelationUserData{UID: followedID}
		currentFollowList = append(currentFollowList, followedUserData)
	} else {
		// If follow is false, remove the followed user from the list if found
		updatedFollowList := make([]*db.RelationUserData, 0, len(currentFollowList)-1)
		for _, user := range currentFollowList {
			if user.UID != followedID {
				updatedFollowList = append(updatedFollowList, user)
			}
		}
		currentFollowList = updatedFollowList
	}

	// Update the follow list in Redis
	updatedFollowListJSON, _ := json.Marshal(currentFollowList)
	err = global.UserInfoRC.Set(followListCacheKey, updatedFollowListJSON, 0).Err()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error updating follow list, ", err.Error())
	}
}

func updateFriendLists(user1ID, user2ID uint64, follow bool) {
	// Update friend list for user1 (user1 follows user2)
	updateFriendList(user1ID, user2ID, follow)

	// Update friend list for user2 (user2 follows user1)
	updateFriendList(user2ID, user1ID, follow)
}

func updateFriendList(userID, friendID uint64, follow bool) {
	// Update the friend list for userID

	friendListCacheKey := fmt.Sprintf("friendList:%d", userID)

	// Get the current friend list from Redis
	currentFriendListJSON, err := global.UserInfoRC.Get(friendListCacheKey).Result()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: ", err.Error())
		return
	}

	var currentFriendList []*api.FriendUser
	if err := json.Unmarshal([]byte(currentFriendListJSON), &currentFriendList); err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error decoding friend list, ", err.Error())
		return
	}

	// Find the friend in the current friend list, or add them if not found
	friendFound := false
	for _, friend := range currentFriendList {
		if uint64(friend.ID) == friendID {
			friendFound = true
			break
		}
	}

	if follow {
		// If follow is true, add the friend if not found in the list
		if !friendFound {
			// TODO:
			userList, _ := db.GetFriendList(userID)
			friendUserList := make([]*api.FriendUser, 0, len(userList))
			for _, u := range userList {
				msg, err := db.GetLatestMsg(userID, u.ID)
				if err != nil {
					hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: ", err.Error())
					return
				}
				friendUserList = append(friendUserList, pack.FriendUser(u, db.IsFollow(userID, u.ID), msg.Content, msg.MsgType))
			}
		}
	} else {
		// If follow is false, remove the friend from the list if found
		if friendFound {
			updatedFriendList := make([]*api.FriendUser, 0, len(currentFriendList)-1)
			for _, friend := range currentFriendList {
				if uint64(friend.ID) != friendID {
					updatedFriendList = append(updatedFriendList, friend)
				}
			}
			currentFriendList = updatedFriendList
		}
	}

	// Update the friend list in Redis
	updatedFriendListJSON, _ := json.Marshal(currentFriendList)
	err = global.UserInfoRC.Set(friendListCacheKey, updatedFriendListJSON, 0).Err()
	if err != nil {
		hlog.Error("pkg.initialize.redis.SubscribeToTopicChanges.relation err: Error updating friend list, ", err.Error())
	}
}

func WarmUpCacheForTopFollowers() error {
	// Retrieve the top-100 users with the most followers
	topFollowers, err := db.SelectTopFollowers(200)
	if err != nil {
		return err
	}

	for _, follower := range topFollowers {
		_, err := service.GetUserInfo(0, follower.ID)
		if err != nil {
			hlog.Error("pkg.initialize.redis.WarmUpCacheForTopFollowers: ", err.Error())
		}
	}

	return nil
}

func WarmUpCacheForTopRelationData() error {
	// Retrieve the top-100 users with the most followers
	topFollowers, err := db.SelectTopFollowers(100)
	if err != nil {
		return err
	}

	for _, follower := range topFollowers {
		// Warm up cache for follow list
		_, err := service.GetFollowList(follower.ID, 0)
		if err != nil {
			hlog.Error("pkg.initialize.redis.WarmUpCacheForTopRelationData.GetFollowList: ", err.Error())
		}

		// Warm up cache for follower list
		_, err = service.GetFollowerList(follower.ID, 0)
		if err != nil {
			hlog.Error("pkg.initialize.redis.WarmUpCacheForTopRelationData.GetFollowerList: ", err.Error())
		}

		// Warm up cache for friend list
		_, err = service.GetFriendList(follower.ID)
		if err != nil {
			hlog.Error("pkg.initialize.redis.WarmUpCacheForTopRelationData.GetFriendList: ", err.Error())
		}
	}

	return nil
}

func WarmUpCacheForTopFavoriteVideos() error {
	// Retrieve the top-500 videos with the most likes
	topVideos, err := db.SelectTopFavoriteVideos(500)
	if err != nil {
		return err
	}

	for _, video := range topVideos {
		// Warm up the cache for this video's like count
		var builder strings.Builder
		builder.WriteString(strconv.FormatUint(video.ID, 10))
		builder.WriteString("_video_like")
		videoLikeKey := builder.String()

		// Check if video like count is available in Redis cache
		_, err := global.VideoFRC.Get(videoLikeKey).Result()
		if err != nil {
			// Cache miss, query the database
			likeInt64, err := db.SelectVideoFavoriteCountByVideoID(video.ID)
			if err != nil {
				hlog.Error("service.favorite.warmUpCacheForVideoLikes err:", err.Error())
				continue
			}

			// Store video like count in Redis cache
			global.VideoFRC.Set(videoLikeKey, likeInt64, 0)
		}
	}

	return nil
}

func WarmUpCacheForTopComments() error {
	// Retrieve the top-500 videos for which you want to warm up the comment cache
	topVideos, err := db.SelectTopFavoriteVideos(500)
	if err != nil {
		return err
	}

	for _, video := range topVideos {
		// Warm up the cache for comments on each video
		// Retrieve comments for the video from the database
		comments, err := db.SelectCommentDataByVideoID(video.ID)
		if err != nil {
			return err
		}

		// Convert comments to the appropriate response format
		commentData := make([]*db.CommentData, len(comments))
		for i, comment := range comments {
			commentData[i] = &db.CommentData{
				CID:         comment.ID,
				Content:     comment.Content,
				CreatedTime: comment.CreatedTime,
				UID:         0,
				Username:    "",
				IsFollow:    false,
				Avatar:      "",
			}
		}

		response := &api.DouyinCommentListResponse{
			StatusCode:  0,
			CommentList: pack.CommentDataList(commentData),
		}

		// Store comments in the Redis cache
		cacheKey := fmt.Sprintf("commentList:%d:%d", 0, video.ID)
		responseJSON, _ := json.Marshal(response)

		// Set cache expiration with some randomization to prevent cache stampede
		cacheDuration := 10*time.Hour + time.Duration(rand.Intn(10))*time.Hour

		err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
