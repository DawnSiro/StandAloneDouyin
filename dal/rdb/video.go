package rdb

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

// VideoInfo 视频固定的信息
type VideoInfo struct {
	ID          uint64
	PublishTime float64
	AuthorID    uint64
	PlayURL     string
	CoverURL    string
	Title       string
}

type VideoInfoCount struct {
	FavoriteCount int64
	CommentCount  int64
}

// FavoriteVideoIDZSet 用户点赞视频ID 集合
type FavoriteVideoIDZSet struct {
	VideoID     uint64
	CreatedTime float64
}

// PublishVideoIDZSet 用户发布视频ID 集合
type PublishVideoIDZSet struct {
	VideoID     uint64
	CreatedTime float64
}

func SetVideoInfo(video *model.Video) error {
	// 拆分成两个数据结构，不做更改或更改频率较低放在 UserInfo 中
	videoInfo := &VideoInfo{
		ID:          video.ID,
		PublishTime: float64(video.PublishTime.UnixMilli()),
		AuthorID:    video.AuthorID,
		PlayURL:     video.PlayURL,
		CoverURL:    video.CoverURL,
		Title:       video.Title,
	}

	// 进行序列化
	videoInfoJSON, err := json.Marshal(videoInfo)
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}

	// 开启管道
	pipeline := global.UserRC.Pipeline()

	// 设置 UserInfo 的 JSON 缓存
	infoKey := constant.VideoInfoRedisPrefix + strconv.FormatUint(video.ID, 10)
	err = pipeline.Set(infoKey, videoInfoJSON,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}

	videoInfoCount := &VideoInfoCount{
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
	}

	infoCountKey := constant.VideoInfoCountHashRedisPrefix + strconv.FormatUint(video.ID, 10)
	// 使用 MSet 进行批量设置
	err = pipeline.HMSet(infoCountKey, map[string]interface{}{
		constant.VideoFavoriteCountRedisFiled: videoInfoCount.FavoriteCount,
		constant.VideoCommentCountRedisFiled:  videoInfoCount.CommentCount,
	}).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	err = pipeline.Expire(infoCountKey,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	// 执行管道中的命令
	_, err = pipeline.Exec()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	return nil
}

func GetVideoInfo(videoID uint64) (*model.Video, error) {
	logTag := "dal.rdb.video.GetVideoInfo err:"
	// 获取用户信息中基本信息，不做更改或更改频率较低的 UserInfo
	infoKey := constant.VideoInfoRedisPrefix + strconv.FormatUint(videoID, 10)
	videoInfoJSON, err := global.UserRC.Get(infoKey).Result()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	videoInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(videoInfoJSON), videoInfo)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	// 获取用户信息中的计数信息，需要频繁更新的 UserInfoCount
	infoCountKey := constant.VideoInfoCountHashRedisPrefix + strconv.FormatUint(videoID, 10)
	countMap, err := global.UserRC.HGetAll(infoCountKey).Result()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	// 进行参数解析
	favoriteCount, err := strconv.ParseInt(countMap[constant.VideoFavoriteCountRedisFiled], 10, 64)
	commentCount, err := strconv.ParseInt(countMap[constant.VideoCommentCountRedisFiled], 10, 64)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}
	videoInfoCount := &VideoInfoCount{
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
	}

	// 更新缓存时间
	err = global.VideoRC.Expire(infoKey,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
	}
	err = global.VideoRC.Expire(infoCountKey,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
	}

	return &model.Video{
		ID:            videoInfo.ID,
		PublishTime:   time.UnixMicro(int64(videoInfo.PublishTime)),
		AuthorID:      videoInfo.AuthorID,
		PlayURL:       videoInfo.PlayURL,
		CoverURL:      videoInfo.CoverURL,
		FavoriteCount: videoInfoCount.FavoriteCount,
		CommentCount:  videoInfoCount.CommentCount,
		Title:         videoInfo.Title,
	}, nil
}

// SetFavoriteVideoID 加载用户点赞视频列表到 Redis
func SetFavoriteVideoID(userID uint64, favoriteVideoIDZSet []*FavoriteVideoIDZSet) error {
	key := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	zList := make([]redis.Z, len(favoriteVideoIDZSet))
	for i, set := range favoriteVideoIDZSet {
		zList[i] = redis.Z{
			Score:  set.CreatedTime,
			Member: set.VideoID,
		}
	}
	err := global.VideoRC.ZAdd(key, zList...).Err()
	if err != nil {
		hlog.Error("rdb.video.FavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

func GetFavoriteVideoID(userID uint64) ([]uint64, error) {
	key := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	strList, err := global.VideoRC.ZRevRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	result := make([]uint64, len(strList))
	for i, str := range strList {
		result[i], err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// AddFavoriteVideoID 添加 Redis 用户点赞视频 ID
func AddFavoriteVideoID(userID uint64, favoriteVideoID *FavoriteVideoIDZSet) error {
	key := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	err := global.VideoRC.ZAdd(key).Err()
	if err != nil {
		hlog.Error("rdb.video.FavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

// DelFavoriteVideoID 删除 Redis 用户点赞视频 ID
func DelFavoriteVideoID(userID, videoID uint64) error {
	key := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	err := global.VideoRC.SRem(key, videoID).Err()
	if err != nil {
		hlog.Error("rdb.video.DelFavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

// FavoriteVideo 点赞视频
// 1. 新增点赞视频ID
// 2. 增加用户点赞数
// 3. 增加作者被点赞数
func FavoriteVideo(userID, authorID uint64, fVideoIDZSet *FavoriteVideoIDZSet) error {
	logTag := "dal.rdb.video.FavoriteVideo err:"
	videoIDKey := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	userInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(userID, 10)
	authorInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(authorID, 10)
	// 逻辑先留这
	//exists, err := global.VideoRC.Exists(videoIDKey).Result()
	//if exists > 0 {
	//	global.VideoRC.ZAdd(videoIDKey, redis.Z{
	//		Score:  fVideoIDZSet.CreatedTime,
	//		Member: fVideoIDZSet.VideoID,
	//	})
	//}
	//exists, err = global.VideoRC.Exists(userInfoCountKey).Result()
	//if exists > 0 {
	//	global.VideoRC.HIncrBy(userInfoCountKey, constant.FavoriteCountRedisFiled, 1)
	//}
	//exists, err = global.VideoRC.Exists(authorInfoCountKey).Result()
	//if exists > 0 {
	//	global.VideoRC.HIncrBy(authorInfoCountKey, constant.TotalFavoritedRedisFiled, 1)
	//}

	keys := []string{videoIDKey, userInfoCountKey, authorInfoCountKey}
	args := []interface{}{fVideoIDZSet.CreatedTime, fVideoIDZSet.VideoID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err := global.UserRC.EvalSha(global.FavoriteVideoLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// CancelFavoriteVideo 取消点赞视频
// 1. 删去点赞视频ID
// 2. 减少用户点赞数
// 3. 减少作者被点赞数
func CancelFavoriteVideo(userID, authorID, videoID uint64) error {
	logTag := "dal.rdb.video.CancelFavorite err:"
	videoIDKey := constant.FavoriteVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	userInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(userID, 10)
	authorInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(authorID, 10)

	keys := []string{videoIDKey, userInfoCountKey, authorInfoCountKey}
	args := []interface{}{videoID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err := global.VideoRC.EvalSha(global.FavoriteVideoLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// PublishVideo 发布视频
// 1. 新增作品视频ID
// 2. 增加自己的作品数
func PublishVideo(userID uint64, pVideoIDZSet *PublishVideoIDZSet) error {
	logTag := "dal.rdb.video.PublishVideo err:"
	videoIDKey := constant.PublishVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	userInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(userID, 10)

	keys := []string{videoIDKey, userInfoCountKey}
	args := []interface{}{pVideoIDZSet.CreatedTime, pVideoIDZSet.VideoID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err := global.VideoRC.EvalSha(global.PublishVideoLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// SetPublishVideoID 加载用户视频列表到 Redis
func SetPublishVideoID(userID uint64, publishVideoIDZSet []*PublishVideoIDZSet) error {
	key := constant.PublishVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	zList := make([]redis.Z, len(publishVideoIDZSet))
	for i, set := range publishVideoIDZSet {
		zList[i] = redis.Z{
			Score:  set.CreatedTime,
			Member: set.VideoID,
		}
	}
	err := global.VideoRC.ZAdd(key, zList...).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetPublishVideoID(userID uint64) ([]uint64, error) {
	key := constant.PublishVideoIDRedisZSetPrefix + strconv.FormatUint(userID, 10)
	strList, err := global.VideoRC.ZRevRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	result := make([]uint64, len(strList))
	for i, str := range strList {
		result[i], err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
