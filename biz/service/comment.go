package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/pulsar"
	"douyin/pkg/util"
	"douyin/pkg/util/sensitive"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Global variables for cache status flag and Bloom filter
var (
	cacheMutex  sync.Mutex
	cacheStatus map[string]*sync.WaitGroup
	bloomFilter *bloom.BloomFilter
)

func init() {
	cacheStatus = make(map[string]*sync.WaitGroup)
	bloomFilter = bloom.NewWithEstimates(1000000, 0.01)
}

func PostComment(userID, videoID uint64, commentText string) (*api.DouyinCommentActionResponse, error) {
	// 删除redis评论列表缓存
	// 使用 strings.Builder 来优化字符串的拼接
	//var builder strings.Builder
	//builder.WriteString(strconv.FormatUint(videoID, 10))
	//builder.WriteString("_video_comments")
	//delCommentListKey := builder.String()
	//hlog.Info("service.comment.PostComment delCommentListKey:", delCommentListKey)

	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errno.ContainsProhibitedSensitiveWordsError
	}

	// 基于雪花算法生成comment_id
	id, err := util.GetSonyflakeID()
	if err != nil {
		hlog.Error("service.comment.PostComment err: failed to create comment id, ", err.Error())
	}

	// 发布消息队列
	msg := pulsar.PostCommentMessage{
		ID:          id,
		VideoID:     videoID,
		UserID:      userID,
		Content:     commentText,
		CreatedTime: time.Now(),
	}
	err = pulsar.GetPostCommentMQInstance().PostComment(msg)
	if err != nil {
		hlog.Error("service.comment.PostComment err: failed to publish mq ", err.Error())
		return nil, err
	}

	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		return nil, err
	}

	// Notify cache invalidation
	// Publish a message to the Redis channel indicating a comment list change
	global.UserInfoRC.Publish("commentList_changes", "comment_added"+"&"+strconv.FormatUint(userID, 10)+"&"+strconv.FormatUint(videoID, 10))

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment((*db.Comment)(&msg), dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	// 查询此评论是否是本人发送的
	isComment := db.IsCommentCreatedByMyself(userID, commentID)
	// 非本人评论
	if !isComment {
		hlog.Error("service.comment.DeleteComment err:", errno.DeletePermissionError)
		return nil, errno.DeletePermissionError
	}

	dbc, err := db.DeleteCommentByID(videoID, commentID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	dbu, err := db.SelectUserByID(userID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}
	authorID, err := db.SelectAuthorIDByVideoID(videoID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}

	// Notify cache invalidation
	// Publish a message to the Redis channel indicating a comment list change
	global.UserInfoRC.Publish("commentList_changes", "comment_deleted"+"&"+strconv.FormatUint(userID, 10)+"&"+strconv.FormatUint(videoID, 10))

	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, db.IsFollow(userID, authorID)),
	}, nil
}

func GetCommentList(userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	// Check if comment data is available in Redis cache
	cacheKey := fmt.Sprintf("commentList:%d:%d", userID, videoID)

	// 解决缓存穿透 -- 添加布隆过滤器判断值是否在于RC或DB
	if bloomFilter.TestString(cacheKey) {
		cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
		if err == nil {
			// Cache hit, return cached comment data
			var cachedResponse api.DouyinCommentListResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
				hlog.Error("service.comment.GetCommentList err: Error decoding cached data, ", err.Error())
			} else {
				return &cachedResponse, nil
			}
		}
	}

	// 解决缓存击穿 -- 添加互斥锁，同一时间仅有一个线程更新缓存
	// Create a WaitGroup for the cache update operation
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Check if another thread is updating the cache
	cacheMutex.Lock()
	existingWaitGroup, exists := cacheStatus[cacheKey]
	if exists {
		cacheMutex.Unlock()
		existingWaitGroup.Wait()
		return GetCommentList(userID, videoID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	commentData, err := db.SelectCommentDataByVideoIDANDUserID(videoID, userID)
	if err != nil {
		hlog.Error("service.comment.GetCommentList err:", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: pack.CommentDataList(commentData),
	}

	// Store comment data in Redis cache
	responseJSON, _ := json.Marshal(response)

	// 解决缓存雪崩 -- 添加随机数，避免缓存同时过期
	// Add a random duration to the cache valid time
	cacheDuration := 10*time.Minute + time.Duration(rand.Intn(600))*time.Second

	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.comment.GetCommentList err: Error storing data in cache, ", err.Error())
	}

	// Add cacheKey to Bloom filter
	bloomFilter.AddString(cacheKey)

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}
