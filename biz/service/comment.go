package service

import (
	"context"
	"douyin/dal/model"
	"douyin/pkg/constant"
	"strconv"
	"sync"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/pulsar"
	"douyin/pkg/util"
	"douyin/pkg/util/sensitive"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	//检测是否带有敏感词
	if sensitive.IsWordsFilter(commentText) {
		return nil, errno.ContainsProhibitedSensitiveWordsError
	}

	// 基于雪花算法生成comment_id
	id, err := util.GetSonyFlakeID()
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

	// 查询缓存数据
	dbu, err := rdb.GetUserInfo(userID)
	if err != nil {
		hlog.Error("service.comment.PostComment err:", err.Error())
		// 缓存没有再查数据库
		dbu, err = db.SelectUserByID(userID)
		if err != nil {
			hlog.Error("service.comment.PostComment err:", err.Error())
			return nil, err
		}
		// 然后设置缓存
		err = rdb.SetUserInfo(dbu)
		if err != nil {
			// 要是设置出错，也不返回，继续执行逻辑
			hlog.Error("service.comment.PostComment err:", err.Error())
		}
	}

	// 这里的 isFollow 直接返回 false ，因为评论人自己当然不能关注自己
	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment((*model.Comment)(&msg), dbu, false),
	}, nil
}

func DeleteComment(userID, videoID, commentID uint64) (*api.DouyinCommentActionResponse, error) {
	// 查询此评论是否是本人发送的
	isComment, err := rdb.IsCommentCreatedByMyself(userID, videoID)
	// 这里因为使用了 string 存储，所以逻辑没有 ZSet 那么复杂
	if err != nil {
		isComment = db.IsCommentCreatedByMyself(userID, commentID)
	}

	// 非本人评论直接返回
	if !isComment {
		hlog.Error("service.comment.DeleteComment err:", errno.DeletePermissionError)
		return nil, errno.DeletePermissionError
	}

	// db 中会一起删除缓存数据
	dbc, err := db.DeleteCommentByID(videoID, commentID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		return nil, err
	}

	// 查询用户评论数据
	dbu, err := rdb.GetUserInfo(userID)
	if err != nil {
		hlog.Error("service.comment.DeleteComment err:", err.Error())
		dbu, err = db.SelectUserByID(userID)
		if err != nil {
			hlog.Error("service.comment.DeleteComment err:", err.Error())
			return nil, err
		}
	}

	// 这里的 isFollow 直接返回 false ，因为评论人自己当然不能关注自己
	return &api.DouyinCommentActionResponse{
		StatusCode: 0,
		Comment:    pack.Comment(dbc, dbu, false),
	}, nil
}

func GetCommentList(ctx context.Context, userID, videoID uint64) (*api.DouyinCommentListResponse, error) {
	// userID可能为0，因为可能存在不登录也能查看视频评论的需求，但是videoID一定得为真实存在的ID
	// 使用布隆过滤器判断视频ID是否存在
	if !global.VideoIDBloomFilter.TestString(strconv.FormatUint(videoID, 10)) {
		hlog.Error("service.comment.GetCommentList err: 布隆过滤器拦截")
		return nil, errno.UserRequestParameterError
	}

	// 加分布式锁
	lock := rdb.NewUserKeyLock(userID, constant.CommentRedisZSetPrefix)
	// 如果 redis 不可用，应该使用程序代码之外的方式进行限流
	_ = lock.Lock(ctx, global.CommentRC)

	// 获取评论基本数据

	// 缓存不存在，查询数据库
	//db.SelectCommentDataByVideoID()

	// 设置缓存

	// 根据用户关注缓存判断用户是否关注了评论的用户

	// 缓存不存在也查询数据库并缓存

	//

	// 获取评论ID
	//commentIDList, err := rdb.GetCommentIDByVideoID(videoID)
	//if err == nil {
	//	dbCList, err := rdb.GetCommentList(commentIDList)
	//	if err == nil {
	//		// 先从缓存中查询用户数据，并把未命中的 UserID 放入切片中
	//		unHitUserID := make([]uint64, 0)
	//		for i := 0; i < len(commentIDList); i++ {
	//			info, err := rdb.GetUserInfo(dbCList[i].UserID)
	//			if err != nil {
	//				unHitUserID = append(unHitUserID, dbCList[i].UserID)
	//			}
	//		}
	//	}
	//
	//}

	// Cache miss, query the database
	commentData, err := db.SelectCommentDataByVideoIDAndUserID(videoID, userID)
	if err != nil {
		hlog.Error("service.comment.GetCommentList err:", err.Error())
		// Release cache status flag to allow other threads to update cache
		//cacheMutex.Lock()
		//delete(cacheStatus, cacheKey)
		//cacheMutex.Unlock()
		return nil, err
	}

	// Convert to response format
	response := &api.DouyinCommentListResponse{
		StatusCode:  0,
		CommentList: pack.CommentDataList(commentData),
	}

	//// Store comment data in Redis cache
	//responseJSON, _ := json.Marshal(response)
	//
	//// 解决缓存雪崩 -- 添加随机数，避免缓存同时过期
	//// Add a random duration to the cache valid time
	//cacheDuration := 10*time.Minute + time.Duration(rand.Intn(600))*time.Second

	//err = global.UserRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	//if err != nil {
	//	hlog.Error("service.comment.GetCommentList err: Error storing data in cache, ", err.Error())
	//}
	//
	//// Add cacheKey to Bloom filter
	//bloomFilter.AddString(cacheKey)
	//
	//// Release cache status flag and signal that cache update is done
	//cacheMutex.Lock()
	//delete(cacheStatus, cacheKey)
	//waitGroup.Done()
	//cacheMutex.Unlock()

	return response, nil
}
